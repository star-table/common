package mq

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/library/mq/dbmq"
	"github.com/star-table/common/library/mq/kafka"
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
