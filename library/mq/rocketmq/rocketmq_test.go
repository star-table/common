package rocketmq

import (
	"fmt"
	rmq "github.com/apache/rocketmq-client-go/core"
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/polaris-backend/polaris-manager/manager/model/bo"
	"sync/atomic"
	"testing"
	"time"
)

func sendMessage(config *rmq.ProducerConfig) {
	producer, err := rmq.NewProducer(config)

	if err != nil {
		fmt.Println("create Producer failed, error:", err)
		return
	}

	err = producer.Start()
	if err != nil {
		fmt.Println("start producer error", err)
		return
	}
	defer producer.Shutdown()

	fmt.Printf("Producer: %s started... \n", producer)
	for i := 2; i < 4; i++ {
		msg := fmt.Sprintf("%s-%d", "msg", i)
		result, err := producer.SendMessageSync(&rmq.Message{Topic: "test_topic", Body: msg})
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Printf("send message: %s result: %s\n", msg, result)
	}
	fmt.Println("shutdown producer.")
}

func pullMessage(config *rmq.PullConsumerConfig) {
	consumer, err := rmq.NewPullConsumer(config)
	if err != nil {
		fmt.Printf("new pull consumer error:%s\n", err)
		return
	}

	err = consumer.Start()
	if err != nil {
		fmt.Printf("start consumer error:%s\n", err)
		return
	}
	defer consumer.Shutdown()
	mqs := consumer.FetchSubscriptionMessageQueues("test_topic")
	fmt.Printf("fetch subscription mqs:%+v\n", mqs)

	total, offsets, now := 0, map[int]int64{}, time.Now()

	fmt.Println(total)
	fmt.Println(offsets)
	fmt.Println(now)

PULL:
	for {
		for _, mq := range mqs {
			pr := consumer.Pull(mq, "*", offsets[mq.ID], 32)
			total += len(pr.Messages)
			fmt.Printf("pull %s, result:%+v\n", mq.String(), pr)

			switch pr.Status {
			case rmq.PullNoNewMsg:
				break PULL
			case rmq.PullFound:
				fallthrough
			case rmq.PullNoMatchedMsg:
				fallthrough
			case rmq.PullOffsetIllegal:
				offsets[mq.ID] = pr.NextBeginOffset
			case rmq.PullBrokerTimeout:
				fmt.Println("broker timeout occur")
			}
		}
	}

	var timePerMessage time.Duration
	if total > 0 {
		timePerMessage = time.Since(now) / time.Duration(total)
	}
	fmt.Printf("total message:%d, per message time:%d\n", total, timePerMessage)

}

func consumeWithPush(config *rmq.PushConsumerConfig) {

	consumer, err := rmq.NewPushConsumer(config)
	if err != nil {
		println("create Consumer failed, error:", err)
		return
	}

	ch := make(chan interface{})
	var count = (int64)(64)
	// MUST subscribe topic before consumer started.
	consumer.Subscribe("test_topic", "*", func(msg *rmq.MessageExt) rmq.ConsumeStatus {
		fmt.Printf("A message received: \"%s\" \n", msg.Body)
		fmt.Printf("A message received: \"%v\" \n", msg)
		if atomic.AddInt64(&count, -1) <= 0 {
			ch <- "quit"
		}
		return rmq.ReConsumeLater
	})

	err = consumer.Start()
	if err != nil {
		println("consumer start failed,", err)
		return
	}

	fmt.Printf("consumer: %s started...\n", consumer)
	<-ch
	err = consumer.Shutdown()
	if err != nil {
		println("consumer shutdown failed")
		return
	}
	println("consumer has shutdown.")
}

func TestGetRocketMQProxy(t *testing.T) {

	config.LoadTreeConfig()

	//pConfig := &rocketmq.ProducerConfig{
	//	ClientConfig: rocketmq.ClientConfig{
	//		GroupID:    config.GetMQ().Rocket.GroupID,
	//		NameServer: config.GetMQ().Rocket.NameServer,
	//		LogC: &rocketmq.LogConfig{
	//			Path:     config.GetMQ().Rocket.Log.LogPath,
	//			FileSize: config.GetMQ().Rocket.Log.FileSize,
	//			FileNum:  config.GetMQ().Rocket.Log.FileNum,
	//			Level:    rocketmq.LogLevelDebug,
	//		},
	//}}
	//
	//if "info" == config.GetMQ().Rocket.Log.Level {
	//	pConfig.LogC.Level = rocketmq.LogLevelInfo
	//}else{
	//	pConfig.LogC.Level = rocketmq.LogLevelDebug
	//}
	//
	//sendMessage(pConfig)

	//cConfig := &rocketmq.PullConsumerConfig{
	//	ClientConfig: rocketmq.ClientConfig{
	//		GroupID:    config.GetMQ().Rocket.GroupID,
	//		NameServer: config.GetMQ().Rocket.NameServer,
	//		LogC: &rocketmq.LogConfig{
	//			Path:     config.GetMQ().Rocket.Log.LogPath,
	//			FileSize: config.GetMQ().Rocket.Log.FileSize,
	//			FileNum:  config.GetMQ().Rocket.Log.FileNum,
	//			Level:    rocketmq.LogLevelDebug,
	//		},
	//}}
	//
	//pullMessage(cConfig)

	//cConfig := &rocketmq.PushConsumerConfig{
	//	ClientConfig: rocketmq.ClientConfig{
	//		GroupID:    config.GetMQ().Rocket.GroupID,
	//		NameServer: config.GetMQ().Rocket.NameServer,
	//		LogC: &rocketmq.LogConfig{
	//			Path:     config.GetMQ().Rocket.Log.LogPath,
	//			FileSize: config.GetMQ().Rocket.Log.FileSize,
	//			FileNum:  config.GetMQ().Rocket.Log.FileNum,
	//			Level:    rocketmq.LogLevelDebug,
	//		},
	//	}}
	//
	//consumeWithPush(cConfig)

	msg := &bo.MqMessage{

		Topic: "test_test_topic",
		Tags:  "a,b,c",
		Keys:  "123",
		Body:  "msg_333",
	}

	msgExt, err := GetRocketMQProxy().SendMessage(msg)

	fmt.Println(msgExt)
	fmt.Println(err)
}

func pushMsg(message *bo.MqMessageExt) errors.SystemErrorInfo {
	fmt.Println(message)
	return nil
}

func TestRocketMQProxy_PushMessage(t *testing.T) {

	config.LoadTreeConfig()

	GetRocketMQProxy().PushMessage("test_test_topic", config.GetMQ().Rocket.GroupID, pushMsg)
}
