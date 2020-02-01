package kafka

import (
	"fmt"
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/errors"
	"github.com/galaxy-book/common/core/model"
	"github.com/galaxy-book/common/core/util/json"
	"strconv"
	"testing"
	"time"
)

const (
	KafkaTestErrorMessage = "error"
)

var topic = "testtest1"
var groupId = "testtest-group"

func TestProxy_SendMessage(t *testing.T) {
	config.LoadUnitTestConfig()
	proxy := Proxy{}

	reconsumer := 5

	for i := 0; i < 10000; i ++{
		_, err := proxy.PushMessage(&model.MqMessage{
			Topic:          topic,
			Body:           strconv.Itoa(i),
			ReconsumeTimes: &reconsumer,
			RePushTimes:    &reconsumer,
		})
		fmt.Println("生产者",err)
		time.Sleep(2 * time.Second)
	}
}

func TestProxy_ConsumePushMessage(t *testing.T) {
	config.LoadUnitTestConfig()

	fmt.Println(json.ToJsonIgnoreError(config.GetMQ().Kafka))

	proxy := Proxy{}
	go proxy.ConsumeMessage(topic, groupId, func(msg *model.MqMessageExt) errors.SystemErrorInfo {
		log.Infof("msg offset: %d, partition: %d,  value: %s", msg.Offset, msg.Partition, string(msg.Body))
		if msg.Body == KafkaTestErrorMessage {
			return errors.BuildSystemErrorInfo(errors.KafkaMqConsumeMsgError)
		}
		return nil
	}, func(message *model.MqMessageExt) {

		log.Info("最终失败:" + json.ToJsonIgnoreError(message))
	})

	select {

	}
}
