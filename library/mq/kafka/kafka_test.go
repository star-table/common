package kafka

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"testing"
)

const(
	KafkaTestErrorMessage = "error"
)

func TestProxy_SendMessage(t *testing.T) {
	config.LoadUnitTestConfig()
	proxy := Proxy{}

	reconsumer := 5
	_, err := proxy.PushMessage(&model.MqMessage{
		Topic:     "unittest",
		Partition: 0,
		Body:      "123",
		ReconsumeTimes: &reconsumer,
		RePushTimes: &reconsumer,
	})
	t.Log(err)
}


func TestProxy_ConsumePushMessage(t *testing.T) {
	config.LoadUnitTestConfig()

	proxy := Proxy{}
	go proxy.ConsumeMessage("unittest", "123", func(msg *model.MqMessageExt) errors.SystemErrorInfo {
		log.Infof("msg offset: %d, partition: %d,  value: %s", msg.Offset, msg.Partition, string(msg.Body))
		if msg.Body == KafkaTestErrorMessage{
			return errors.BuildSystemErrorInfo(errors.KafkaMqConsumeMsgError)
		}
		return nil
	}, func(message *model.MqMessageExt){

		log.Info("最终失败:", json.ToJsonIgnoreError(message))
	})
}

