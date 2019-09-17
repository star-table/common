package kafka

import (
	"context"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/library/mq"
	"github.com/Shopify/sarama"
	"strings"
)

type exampleConsumerGroupHandler struct {
	fu func(message *mq.MqMessageExt) errors.SystemErrorInfo
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

	select {
	case msg := <-claim.Messages():
		if msg != nil {
			msgExt := &mq.MqMessageExt{
				MqMessage: mq.MqMessage{
					Topic:     msg.Topic,
					Body:      string(msg.Value),
					Keys:      string(msg.Key),
					Partition: msg.Partition,
					Offset:    msg.Offset,
				},
			}
			err1 := h.fu(msgExt)
			if err1 != nil {
				log.Errorf("Kafka 业务消费异常 %v", err1)
				log.Errorf("Topic: %s, 处理失败的消息：%s ", msg.Topic, string(msg.Value))
			} else {
				//暂时自动提交，之后考虑嵌入的自定义方法中
				sess.MarkMessage(msg, "")
			}
		}
	}
	//}
	//fmt.Println("ConsumeClaim End", claim)
	return nil
}

func (proxy *Proxy) ConsumeMessage(topic string, groupId string, fu func(message *mq.MqMessageExt) errors.SystemErrorInfo) errors.SystemErrorInfo {
	kafkaConfig := getKafkaConfig()
	log.Infof("Starting a new Sarama consumer, topic %s, groupId %s", topic, groupId)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = version

	topics := []string{topic}
	ctx, _ := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(kafkaConfig.NameServers, ","), groupId, config)

	if err != nil {
		log.Errorf("Error creating consumer group client: %v", err)
		return errors.BuildSystemErrorInfo(errors.KafkaMqConsumeStartError)
	}

	handler := exampleConsumerGroupHandler{
		fu: fu,
	}
	for {
		//log.Infof("准备消费, topic %s, groupId %s", topic, groupId)

		if err := client.Consume(ctx, topics, &handler); err != nil {
			log.Errorf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			log.Errorf("异常退出 %v", ctx.Err())
		}
	}

	log.Info("消费结束")

	return nil
}
