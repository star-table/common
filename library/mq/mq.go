package mq

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/library/mq/dbmq"
	"gitea.bjx.cloud/allstar/common/library/mq/kafka"
)

type MQClient interface {
	PushMessage(messages ...*model.MqMessage) (*[]model.MqMessageExt, errors.SystemErrorInfo)
	ConsumeMessage(topic string, groupId string, fu func(message *model.MqMessageExt) errors.SystemErrorInfo, errCallback func(message *model.MqMessageExt)) errors.SystemErrorInfo
}

var (
	kafkaClient	   MQClient = &kafka.Proxy{}
	dbMqClient		MQClient = &dbmq.DbMQProxy{}
)

func GetMQClient() *MQClient {
	if consts.MQModeKafka == config.GetMQ().Mode {
		return &kafkaClient
	}
	return &dbMqClient
}
