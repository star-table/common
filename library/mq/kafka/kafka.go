package kafka

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"gitea.bjx.cloud/allstar/common/core/util/strs"
	"gitea.bjx.cloud/allstar/common/core/util/uuid"
	"gitea.bjx.cloud/allstar/common/library/cache"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type Proxy struct {
	//key: topic + partitioner
	producers map[string]sarama.AsyncProducer
}

var (
	log     = logger.GetMQLogger()
	version = sarama.V2_3_0_0

	producerConfig = sarama.NewConfig()
)

const (
	RecordHeaderReconsumeTimes = "ReconsumeTimes"
	RecordHeaderRePushTimes    = "RePushTimes"
)

func init() {
	//生产者通用配置
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.Return.Errors = true
	producerConfig.Version = version
}

type MsgMetadata struct {
	//消费重试次数
	ReconsumeTimes int
	//推送重试次数
	RePushTimes int
}

func getKafkaConfig() *config.KafkaMQConfig {
	return config.GetKafkaConfig()
}

func (proxy *Proxy) getProducerAutoConnect(topic string) (*sarama.AsyncProducer, errors.SystemErrorInfo) {
	//key := topic + "#" + strconv.Itoa(int(partition))
	producer, err := proxy.getProducer(topic)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return nil, err
	}
	return producer, nil
}

func (proxy *Proxy) getProducer(topic string) (*sarama.AsyncProducer, errors.SystemErrorInfo) {
	key := topic
	if proxy.producers == nil {
		proxy.producers = map[string]sarama.AsyncProducer{}
	}

	if v, ok := proxy.producers[key]; ok && v != nil{
		return &v, nil
	}

	uuid := uuid.NewUuid()
	suc, err := cache.TryGetDistributedLock(key, uuid)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return nil, errors.BuildSystemErrorInfo(errors.GetDistributedLockError, err)
	}
	if suc {
		//如果获取到锁，则开始初始化
		defer func() {
			if _, err := cache.ReleaseDistributedLock(key, uuid); err != nil{
				log.Error(err)
			}
		}()
	}

	//二次确认
	if v, ok := proxy.producers[key]; ok && v != nil{
		return &v, nil
	}

	//重新构造producer
	producer, err1 := proxy.buildProducer()
	if err1 != nil {
		log.Error(strs.ObjectToString(err1))
		return nil, err1
	}

	proxy.producers[key] = *producer
	return producer, nil
}

func (proxy *Proxy) CloseConnect(topic string) errors.SystemErrorInfo{
	proxy.producers[topic] = nil
	return nil
}

func (proxy *Proxy) buildProducer() (*sarama.AsyncProducer, errors.SystemErrorInfo) {
	kafkaConfig := getKafkaConfig()
	log.Infof("build producer")

	producer, err := sarama.NewAsyncProducer(strings.Split(kafkaConfig.NameServers, ","), producerConfig)
	if err != nil {
		log.Infof("producer_test create producer error :%#v", err)
		return nil, errors.BuildSystemErrorInfo(errors.KafkaMqSendMsgError, err)
	}
	return &producer, nil
}

func (proxy *Proxy) PushMessage(messages ...*model.MqMessage) (*[]model.MqMessageExt, errors.SystemErrorInfo) {
	if messages == nil || len(messages) == 0 {
		return nil, errors.BuildSystemErrorInfo(errors.KafkaMqSendMsgCantBeNullError)
	}

	msgExts := make([]model.MqMessageExt, len(messages))
	for i, message := range messages {
		//传递metadata，方便消费端重试
		ReconsumeTimes := config.GetKafkaConfig().ReconsumeTimes
		RePushTimes := config.GetKafkaConfig().RePushTimes
		if message.ReconsumeTimes != nil {
			ReconsumeTimes = *message.ReconsumeTimes
		}
		if message.RePushTimes != nil {
			RePushTimes = *message.RePushTimes
		}

		key := uuid.NewUuid()
		// send message
		msg := &sarama.ProducerMessage{
			Topic:     message.Topic,
			Partition: message.Partition,
			Key:       sarama.StringEncoder(key),
			Value:     sarama.ByteEncoder(message.Body),
			Headers: []sarama.RecordHeader{
				{
					Key:   []byte(RecordHeaderReconsumeTimes),
					Value: []byte(strconv.Itoa(ReconsumeTimes)),
				},
			},
		}


		var pushErr error = nil
		for rePushTime := 0; rePushTime <= RePushTimes; rePushTime++ {
			p, err1 := proxy.getProducerAutoConnect(message.Topic)
			if err1 != nil {
				log.Error(strs.ObjectToString(err1))
				return nil, err1
			}
			producer := *p

			if rePushTime > 0 {
				log.Infof("重试次数%d，最大次数%d, 上次失败原因%v, 消息内容%s", rePushTime, message.RePushTimes, pushErr, json.ToJsonIgnoreError(message))
			}
			producer.Input() <- msg
			select {
			case suc := <-producer.Successes():
				log.Infof("推送成功, offset: %d,  timestamp: %s， 消息内容%s", suc.Offset, suc.Timestamp.String(), json.ToJsonIgnoreError(message))
				pushErr = nil
			case fail := <-producer.Errors():
				log.Errorf("err: %s\n", fail.Err.Error())
				pushErr = fail.Err
				//return nil, errors.BuildSystemErrorInfo(errors.KafkaMqSendMsgError, fail)

				if pushErr == sarama.ErrNotConnected || pushErr == sarama.ErrClosedClient || pushErr == sarama.ErrOutOfBrokers{
					log.Errorf("断开连接... %v", pushErr)
					//重连
					closeErr := proxy.CloseConnect(message.Topic)
					if closeErr != nil{
						log.Error(closeErr)
					}
				}
				time.Sleep(time.Duration(3) * time.Second)
			}
			if pushErr == nil {
				break
			}
		}
		if pushErr != nil {
			//最终推送失败，记log
			log.Error("消息推送失败，无重试次数", zap.String(consts.LogMqMessageKey, json.ToJsonIgnoreError(message)))
			return nil, errors.BuildSystemErrorInfo(errors.KafkaMqSendMsgError, pushErr)
		}
		log.Info("消息发送成功 %s", zap.String(consts.LogMqMessageKey, json.ToJsonIgnoreError(message)))
		msgExts[i] = model.MqMessageExt{
			MqMessage: model.MqMessage{
				Topic:     msg.Topic,
				Body:      message.Body,
				Keys:      key,
				Partition: msg.Partition,
				Offset:    msg.Offset,
			},
		}
	}
	return &msgExts, nil
}
