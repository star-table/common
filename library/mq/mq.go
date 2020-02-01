package mq

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/consts"
	"github.com/galaxy-book/common/core/errors"
	"github.com/galaxy-book/common/core/model"
	"github.com/galaxy-book/common/library/mq/dbmq"
	"github.com/galaxy-book/common/library/mq/kafka"
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
