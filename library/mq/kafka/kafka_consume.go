package kafka

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/json"
)

type exampleConsumerGroupHandler struct {
	fu          func(message *model.MqMessageExt) errors.SystemErrorInfo
	errCallback func(message *model.MqMessageExt)
	proxy       *Proxy
}

func (exampleConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	return nil
}
func (exampleConsumerGroupHandler) Cleanup(c sarama.ConsumerGroupSession) error {
	return nil
}
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//fix 死循环
	//for {
	//log.Infof("ConsumeClaim - 开始消费 topic %s", claim.Topic())

	for msg := range claim.Messages() {
		if msg != nil {
			//获取重试次数
			ReconsumeTimes := 0
			if msg.Headers != nil {
				for _, header := range msg.Headers {
					if string(header.Key) == RecordHeaderReconsumeTimes {
						v, _ := strconv.ParseInt(string(header.Value), 10, 32)
						ReconsumeTimes = int(v)
					}
				}
			}
			afterConsumerTimes := ReconsumeTimes - 1

			msgExt := &model.MqMessageExt{
				MqMessage: model.MqMessage{
					Topic:          msg.Topic,
					Body:           string(msg.Value),
					Keys:           string(msg.Key),
					Partition:      msg.Partition,
					Offset:         msg.Offset,
					ReconsumeTimes: &afterConsumerTimes,
				},
			}
			err1 := h.fu(msgExt)
			if err1 != nil {
				log.Errorf("Kafka 业务消费异常 %v", err1)
				log.Errorf("Topic: %s, 处理失败的消息：%q", msg.Topic, json.ToJsonBytesIgnoreError(msgExt.MqMessage))
				sess.MarkMessage(msg, "consumer err")

				log.Infof("剩余消费次数%d", afterConsumerTimes)
				if afterConsumerTimes > -1 {
					_, pushErr := h.proxy.PushMessage(&msgExt.MqMessage)
					if pushErr != nil {
						log.Errorf("重试推送失败, 消息内容: %q", json.ToJsonBytesIgnoreError(msgExt.MqMessage))
					}
				} else {
					log.Errorf("无重试次数, 消息最终消费失败, 消息内容: %q", json.ToJsonBytesIgnoreError(msgExt.MqMessage))
					h.errCallback(msgExt)
				}
			} else {
				sess.MarkMessage(msg, "")
			}
		}
	}
	return nil
}

func (proxy *Proxy) ConsumeMessage(topic string, groupId string, fu func(message *model.MqMessageExt) errors.SystemErrorInfo, errCallback func(message *model.MqMessageExt)) errors.SystemErrorInfo {
	kafkaConfig := getKafkaConfig()
	log.Infof("Kafka config %q", json.ToJsonBytesIgnoreError(kafkaConfig))
	log.Infof("Starting a new Sarama consumer, topic %s, groupId %s", topic, groupId)

	if len(topic) == 0 {
		panic("consume on a empty topic!")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	//config.ChannelBufferSize = 2560
	//config.Consumer.Fetch.Min = 1024 * 1024
	//config.Consumer.Fetch.Default = 1024 * 1024 * 2
	//config.Consumer.Fetch.Max = 1024 * 1024 * 10
	//config.Consumer.MaxWaitTime = 2 * time.Second
	//config.Consumer.MaxProcessingTime = 1 * time.Second
	//config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Version = version

	topics := []string{topic}

	for {

		log.Info("开始连接...")

		ctx, _ := context.WithCancel(context.Background())
		client, err := sarama.NewConsumerGroup(strings.Split(kafkaConfig.NameServers, ","), groupId, config)
		if err != nil {
			log.Errorf("Error creating consumer group client: %v", err)
			return errors.BuildSystemErrorInfo(errors.KafkaMqConsumeStartError)
		}

		handler := exampleConsumerGroupHandler{
			fu:          fu,
			proxy:       proxy,
			errCallback: errCallback,
		}

		for {
			//log.Infof("准备消费, topic %s, groupId %s", topic, groupId)
			if err := client.Consume(ctx, topics, &handler); err != nil {
				log.Errorf("[%s] Error from consumer: %v", topic, err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Errorf("异常退出 %v", ctx.Err())
				break
			}
		}

		err = client.Close()
		if err != nil {
			log.Errorf("关闭连接失败 %v", err)
		}
		time.Sleep(2 * time.Second)
		log.Info("准备重连...")
	}

	log.Info("消费结束")

	return nil
}
