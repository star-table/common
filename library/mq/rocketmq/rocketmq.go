package rocketmq

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/core/util/copyer"
	rocketmq "github.com/apache/rocketmq-client-go/core"
	"sync"
	"sync/atomic"
)

var lock sync.Mutex
var log = logger.GetDefaultLogger()

type Proxy struct {
	pConfig  *rocketmq.ProducerConfig
	producer *rocketmq.Producer

	cConfig  *rocketmq.PushConsumerConfig
	consumer map[string]rocketmq.PushConsumer
}

func (rocketMQProxy *Proxy) initProducer() errors.SystemErrorInfo {
	if rocketMQProxy.pConfig == nil || rocketMQProxy.producer == nil {
		lock.Lock()
		defer lock.Unlock()

		if rocketMQProxy.pConfig == nil {
			rocketMQProxy.pConfig = &rocketmq.ProducerConfig{
				ClientConfig: rocketmq.ClientConfig{
					GroupID:    config.GetMQ().Rocket.GroupID,
					NameServer: config.GetMQ().Rocket.NameServer,
					LogC: &rocketmq.LogConfig{
						Path:     config.GetMQ().Rocket.Log.LogPath,
						FileSize: config.GetMQ().Rocket.Log.FileSize,
						FileNum:  config.GetMQ().Rocket.Log.FileNum,
						Level:    rocketmq.LogLevelInfo,
					},
				}}

			if "debug" == config.GetMQ().Rocket.Log.Level {
				rocketMQProxy.pConfig.LogC.Level = rocketmq.LogLevelDebug
			}
		}

		if rocketMQProxy.producer == nil {
			pro, err2 := rocketmq.NewProducer(rocketMQProxy.pConfig)
			if err2 != nil {
				log.Error("init RocketMQ produce fail. ", err2)
				return errors.BuildSystemErrorInfo(errors.RocketMQProduceInitError, err2)
			}
			rocketMQProxy.producer = &pro
			err2 = (*(rocketMQProxy.producer)).Start()
			if err2 != nil {
				log.Error("init RocketMQ produce fail. ", err2)
				return errors.BuildSystemErrorInfo(errors.RocketMQProduceInitError, err2)
			}
		}
	}
	return nil
}

func (rocketMQProxy *Proxy) initConsumer(groupId string) errors.SystemErrorInfo {
	if rocketMQProxy.pConfig == nil || rocketMQProxy.consumer == nil {
		lock.Lock()
		defer lock.Unlock()

		if rocketMQProxy.pConfig == nil {
			rocketMQProxy.cConfig = &rocketmq.PushConsumerConfig{
				ClientConfig: rocketmq.ClientConfig{
					GroupID:    groupId,
					NameServer: config.GetMQ().Rocket.NameServer,
					LogC: &rocketmq.LogConfig{
						Path:     config.GetMQ().Rocket.Log.LogPath,
						FileSize: config.GetMQ().Rocket.Log.FileSize,
						FileNum:  config.GetMQ().Rocket.Log.FileNum,
						Level:    rocketmq.LogLevelInfo,
					},
				}}

			if "debug" == config.GetMQ().Rocket.Log.Level {
				rocketMQProxy.pConfig.LogC.Level = rocketmq.LogLevelDebug
			}
		}

		if rocketMQProxy.consumer == nil {
			rocketMQProxy.consumer = map[string]rocketmq.PushConsumer{}
		}
	}
	return nil
}

func (rocketMQProxy *Proxy) SendMessage(messages ...*model.MqMessage) (*[]model.MqMessageExt, errors.SystemErrorInfo) {
	err := rocketMQProxy.initProducer()
	if err != nil {
		return nil, err
	}

	msgExts := make([]model.MqMessageExt, 0, len(messages))
	for i, msg := range messages {

		msgExts = append(msgExts, model.MqMessageExt{
			MqMessage: *msg,
		})
		result, err := (*(rocketMQProxy.producer)).SendMessageSync(&rocketmq.Message{Topic: msg.Topic, Body: msg.Body})
		if err != nil {
			log.Error(errors.RocketMQSendMsgError, err)
			msgExts[i].SendStatus = consts.SendMQStatusFail
			continue
		}
		msgExts[i].MessageID = result.MsgId
		msgExts[i].QueueOffset = result.Offset
		if rocketmq.SendOK == result.Status {
			msgExts[i].SendStatus = consts.SendMQStatusSuccess
		} else {
			msgExts[i].SendStatus = consts.SendMQStatusFail
		}
	}

	return &msgExts, nil
}

func (rocketMQProxy *Proxy) ConsumePushMessage(topic string, groupId string, fu func(message *model.MqMessageExt) errors.SystemErrorInfo) errors.SystemErrorInfo {

	err := rocketMQProxy.initConsumer(groupId)
	if err != nil {
		log.Error("create Consumer failed, error:", err)
		return errors.BuildSystemErrorInfo(errors.RocketMQConsumerInitError, err)
	}

	key := topic + "#" + groupId

	if _, ok := rocketMQProxy.consumer[key]; !ok {
		consumer, err := rocketmq.NewPushConsumer(rocketMQProxy.cConfig)
		if err != nil {
			log.Error("create Consumer failed, error:", err)
			return errors.BuildSystemErrorInfo(errors.RocketMQConsumerInitError, err)
		}
		rocketMQProxy.consumer[key] = consumer
	}

	ch := make(chan interface{})
	var count = (int64)(64)

	co := rocketMQProxy.consumer[key]
	// MUST subscribe topic before consumer started.
	co.Subscribe(topic, "*", func(msg *rocketmq.MessageExt) rocketmq.ConsumeStatus {

		msgExt := &model.MqMessageExt{}
		err := copyer.Copy(msg, msgExt)
		if err != nil {
			return rocketmq.ReConsumeLater
		}

		log.Info(" receive msg : %v ", msgExt)
		if atomic.AddInt64(&count, -1) <= 0 {
			ch <- "quit"
		}

		err1 := fu(msgExt)
		if err1 != nil && err1.Code() != errors.OK.Code() {
			return rocketmq.ReConsumeLater
		}
		return rocketmq.ConsumeSuccess
	})

	err2 := co.Start()
	if err2 != nil {
		log.Error("consumer start failed ", err2, topic)
		return errors.BuildSystemErrorInfo(errors.RocketMQConsumerStartError, err2)
	}

	log.Info("consumer: %s topic: %s started...", co, topic)
	<-ch
	err2 = co.Shutdown()
	if err2 != nil {
		log.Error("consumer shutdown failed ", err2, topic)
		return errors.BuildSystemErrorInfo(errors.RocketMQConsumerStopError, err2)
	}

	return nil
}

func GetRocketMQProxy() *Proxy {
	return &Proxy{}
}
