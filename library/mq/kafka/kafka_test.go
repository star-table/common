package kafka

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/model"
	"testing"
	"time"
)

func TestProxy_SendMessage(t *testing.T) {
	config.LoadUnitTestConfig()
	proxy := Proxy{}

	_, err := proxy.PushMessage(&model.MqMessage{
		Topic:     "unittest",
		Partition: 0,
		Body:      "hello",
	})
	t.Log(err)
}


func TestProxy_ConsumePushMessage(t *testing.T) {
	config.LoadUnitTestConfig()

	proxy := Proxy{}
	go proxy.ConsumeMessage("unittest", "123", func(msg *model.MqMessageExt) errors.SystemErrorInfo {
		log.Infof("msg offset: %d, partition: %d,  value: %s", msg.Offset, msg.Partition, string(msg.Body))
		return nil
	})
	time.Sleep(time.Duration(3) * time.Second)
}

