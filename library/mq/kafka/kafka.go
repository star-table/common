package kafka

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"gitea.bjx.cloud/allstar/common/core/util/uuid"
	"gitea.bjx.cloud/allstar/common/library/cache"
	"github.com/Shopify/sarama"
	"strconv"
	"strings"
)

type Proxy struct {
	//key: topic + partitioner
	producers map[string]sarama.AsyncProducer
}

var(
	log = logger.GetMQLogger()
	version = sarama.V2_3_0_0

	producerConfig = sarama.NewConfig()
)

func init(){
	//生产者通用配置
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.Return.Errors = true
	producerConfig.Version = version
}


func getKafkaConfig() *config.KafkaMQConfig{
	return config.GetKafkaConfig()
}

func (proxy *Proxy) getProducerAutoConnect(topic string, partition int32) (*sarama.AsyncProducer, errors.SystemErrorInfo){
	//key := topic + "#" + strconv.Itoa(int(partition))
	producer, err := proxy.getProducer(topic, partition)
	if err != nil{
		log.Error(err)
		return nil, err
	}
	p := *producer

	p.Errors()

	//TODO(nico) 这里做一下断开重连的逻辑
	return producer, nil
}

func (proxy *Proxy) getProducer(topic string, partition int32) (*sarama.AsyncProducer, errors.SystemErrorInfo){
	key := topic + "#" + strconv.Itoa(int(partition))
	if proxy.producers == nil{
		proxy.producers = map[string]sarama.AsyncProducer{}
	}

	if v, ok := proxy.producers[key]; ok {
		return &v, nil
	}

	uuid := uuid.NewUuid()
	suc, err := cache.TryGetDistributedLock(key, uuid)
	if err != nil{
		log.Error(err)
		return nil, errors.BuildSystemErrorInfo(errors.GetDistributedLockError, err)
	}
	if suc{
		//如果获取到锁，则开始初始化
		defer cache.ReleaseDistributedLock(key, uuid)
	}

	//二次确认
	if v, ok := proxy.producers[key]; ok {
		return &v, nil
	}

	//重新构造producer
	producer, err1 := proxy.buildProducer()
	if err1 != nil{
		log.Error(err1)
		return nil, err1
	}

	proxy.producers[key] = *producer
	return producer, nil
}

func (proxy *Proxy) buildProducer() (*sarama.AsyncProducer, errors.SystemErrorInfo){
	kafkaConfig := getKafkaConfig()
	log.Infof("build producer \n")

	producer, err := sarama.NewAsyncProducer(strings.Split(kafkaConfig.NameServers, ","), producerConfig)
	if err != nil {
		log.Infof("producer_test create producer error :%s\n", err.Error())
		return nil, errors.BuildSystemErrorInfo(errors.KafkaMqSendMsgError, err)
	}
	return &producer, nil
}

func (proxy *Proxy) PushMessage(messages ...*model.MqMessage) (*[]model.MqMessageExt, errors.SystemErrorInfo) {
	if messages == nil || len(messages) == 0{
		return nil, errors.BuildSystemErrorInfo(errors.KafkaMqSendMsgCantBeNullError)
	}

	msgExts := make([]model.MqMessageExt, len(messages))
	for i, message := range messages{

		key := uuid.NewUuid()
		// send message
		msg := &sarama.ProducerMessage{
			Topic: message.Topic,
			Partition: message.Partition,
			Key:   sarama.StringEncoder(key),
			Value: sarama.ByteEncoder(message.Body),
		}
		p, err1 := proxy.getProducerAutoConnect(message.Topic, message.Partition)
		if err1 != nil{
			log.Error(err1)
			return nil, err1
		}
		producer := *p

		producer.Input() <- msg
		select {
		case suc := <-producer.Successes():
			log.Infof("offset: %d,  timestamp: %s", suc.Offset, suc.Timestamp.String())
		case fail := <-producer.Errors():
			log.Errorf("err: %s\n", fail.Err.Error())
			return nil, errors.BuildSystemErrorInfo(errors.KafkaMqSendMsgError, fail)
		}
		msgExts[i] = model.MqMessageExt{
			MqMessage: model.MqMessage{
				Topic: msg.Topic,
				Body: message.Body,
				Keys: key,
				Partition: msg.Partition,
				Offset: msg.Offset,
			},
		}
		log.Infof("消息发送成功 %s", json.ToJsonIgnoreError(msgExts))
	}
	return &msgExts, nil
}
