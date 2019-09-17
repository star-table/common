package kafka

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/library/mq"
	"testing"
)

func TestProxy_ConsumePushMessage(t *testing.T) {
	config.LoadLocalConfig()

	proxy := Proxy{}
	_ = proxy.ConsumeMessage("test", "", func(msg *model.MqMessageExt) errors.SystemErrorInfo {
		log.Infof("msg offset: %d, partition: %d,  value: %s", msg.Offset, msg.Partition, string(msg.Body))
		return nil
	})

}

func TestProxy_SendMessage(t *testing.T) {
	config.LoadLocalConfig()
	proxy := Proxy{}

	_, err := proxy.PushMessage(&model.MqMessage{
		Topic:     "test",
		Partition: 0,
		Body:      "hello",
	})
	t.Log(err)
}

func TestProxy_ConsumeGroupMessage(t *testing.T) {
	config.LoadLocalConfig()

	proxy := Proxy{}

	proxy.ConsumeMessage("testabc", "test1", func(msg *model.MqMessageExt) errors.SystemErrorInfo {
		log.Infof("msg offset: %d, partition: %d,  value: %s", msg.Offset, msg.Partition, string(msg.Body))
		return nil
	})

}
