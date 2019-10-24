package dbmq

import (
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/core/util/id/snowflake"
	"gitea.bjx.cloud/allstar/common/core/util/strs"
	"gitea.bjx.cloud/allstar/common/core/util/times"
	"gitea.bjx.cloud/allstar/common/library/db/mysql"
	"strconv"
	"sync"
	"time"
	"upper.io/db.v3"
)

var (
	lock          sync.Mutex
	log           = logger.GetDefaultLogger()
	consumerCount = 10
	sleepTime     = int64(2)
	consumers     = map[string]int{}
	maxRetryCount = 6
)

type DbMQProxy struct{}

func (dbMqProxy *DbMQProxy) PushMessage(messages ...*model.MqMessage) (*[]model.MqMessageExt, errors.SystemErrorInfo) {

	msgCount := len(messages)

	messagePosInterface := make([]interface{}, 0, msgCount)
	for _, v := range messages {
		messagePosInterface = append(messagePosInterface, PpmMqsMessageQueue{
			Id:         snowflake.Id(),
			Topic:      v.Topic,
			MessageKey: v.Keys,
			Message:    v.Body,
			Status:     1,
			Creator:    1,
			CreateTime: time.Now(),
			Updator:    1,
			UpdateTime: time.Now(),
			Version:    1,
			IsDelete:   2,
		})

	}

	err2 := mysql.BatchInsert(&PpmMqsMessageQueue{}, messagePosInterface)
	if err2 != nil {
		log.Error(strs.ObjectToString(err2))
		return nil, errors.BuildSystemErrorInfo(errors.DbMQSendMsgError)
	}

	msgExts := make([]model.MqMessageExt, 0, msgCount)
	for _, v := range messages {
		id := snowflake.Id()
		msgExts = append(msgExts, model.MqMessageExt{
			MqMessage:      *v,
			MessageID:      strconv.FormatInt(id, 10),
			QueueId:        0,
			ReconsumeTimes: 0,
			StoreSize:      0,
			BornTimestamp:  times.GetNowMillisecond(),
			StoreTimestamp: times.GetNowMillisecond(),

			QueueOffset:               id,
			CommitLogOffset:           0,
			PreparedTransactionOffset: 0,
			SendStatus:                1,
		})
	}
	return &msgExts, nil
}

func (dbMqProxy *DbMQProxy) ConsumeMessage(topic string, groupId string, fu func(message *model.MqMessageExt) errors.SystemErrorInfo, errCallback func(message *model.MqMessageExt)) errors.SystemErrorInfo {

	consumer, err := getTopicConsumer(topic, groupId)
	if err != nil {
		return err
	}

	// 需要考虑多实例的同consumer消费的问题，可以用zookeeper解决，目前场景单实例，因此暂不考虑

	for i := 0; i < 1; {
		mqs, err := queryTopicMessage(consumer, consumerCount)

		if checkComsumerTopicMessage(mqs, err) {
			continue
		}

		mqcfm, err := queryTopicGroupMessageFail(consumer, mqs)
		if err != nil {
			log.Error(strs.ObjectToString(err))
			times.Sleep(sleepTime)
			continue
		}

		consumerMessageId, consumerFailMessageId := consumerHandler(consumer, mqs, mqcfm, fu)

		//重试消费mq 和 确认消费
		if retryConsumer(consumer, consumerMessageId, consumerFailMessageId, mqcfm) {
			continue
		}

	}

	return nil
}

//校验topic中的消息
func checkComsumerTopicMessage(mqs *[]PpmMqsMessageQueue, err errors.SystemErrorInfo) (isContinue bool) {

	if err != nil {
		log.Error(strs.ObjectToString(err))
		times.Sleep(sleepTime)
		return true
	}
	if len(*mqs) < 1 {
		times.Sleep(sleepTime)
		return true
	}

	return false
}

func retryConsumer(consumer *PpmMqsMessageQueueConsumer, consumerMessageId int64, consumerFailMessageId int64,

	mqcfm *map[int64]PpmMqsMessageQueueConsumerFail) (isContinue bool) {

	if consumerFailMessageId > 0 {
		b, err := consumerFailMessageHandler(consumer, consumerFailMessageId, mqcfm)
		if err != nil {
			times.Sleep(sleepTime)
			return true
		}
		if b == true {
			// 达到最大重试次数，忽略
			consumerMessageId = consumerFailMessageId
		}
	}

	if consumerMessageId > 0 {
		//有成功消费，更新id
		err := consumerSuccessMessageHandler(consumer, consumerMessageId)
		if err != nil {
			times.Sleep(sleepTime)
			return true
		}
	}

	return false
}

func consumerSuccessMessageHandler(consumer *PpmMqsMessageQueueConsumer, successId int64) errors.SystemErrorInfo {

	consumer.MessageId = successId
	consumer.LastConsumerTime = time.Now()
	consumer.UpdateTime = time.Now()
	err := mysql.Update(consumer)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}
	return nil

}

func consumerFailMessageHandler(consumer *PpmMqsMessageQueueConsumer, failId int64, mqcfm *map[int64]PpmMqsMessageQueueConsumerFail) (bool, errors.SystemErrorInfo) {
	//有成功消费，更新id
	b := false
	if mqcf, ok := (*mqcfm)[failId]; ok {
		// 已有失败记录
		mqcf.FailCount += 1
		mqcf.FailTime = time.Now()
		mqcf.UpdateTime = time.Now()
		if mqcf.FailCount >= maxRetryCount {
			mqcf.Status = consts.MQStatusFail
			b = true
		} else {
			mqcf.Status = consts.MQStatusWait
		}

		err := mysql.Update(&mqcf)
		if err != nil {
			log.Error(strs.ObjectToString(err))
			return b, errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
		}
	} else {
		// 首次失败
		fid := snowflake.Id()

		cf := &PpmMqsMessageQueueConsumerFail{
			Id:         fid,
			Topic:      consumer.Topic,
			GroupName:  consumer.GroupName,
			MessageId:  failId,
			FailCount:  1,
			FailTime:   time.Now(),
			Status:     consts.MQStatusWait,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			Version:    1,
			IsDelete:   2,
		}
		err2 := mysql.Insert(cf)
		if err2 != nil {
			log.Errorf("consumer fail save db fail. %v %v.", *cf, err2)
			return b, errors.BuildSystemErrorInfo(errors.MysqlOperateError, err2)
		}
	}
	return b, nil
}

func consumerHandler(consumer *PpmMqsMessageQueueConsumer, mqs *[]PpmMqsMessageQueue, mqcfm *map[int64]PpmMqsMessageQueueConsumerFail,
	fu func(message *model.MqMessageExt) errors.SystemErrorInfo) (int64, int64) {

	consumerMessageId := int64(0)
	consumerFailMessageId := int64(0)
	for _, v := range *mqs {

		fcount := 0
		if cf, ok := (*mqcfm)[v.Id]; ok {
			fcount = cf.FailCount
		}

		ts := times.GetMillisecond(v.CreateTime)
		msgExt := &model.MqMessageExt{
			MqMessage: model.MqMessage{
				Topic:          consumer.Topic,
				Keys:           v.MessageKey,
				Body:           v.Message,
				DelayTimeLevel: fcount,
			},
			MessageID:       strconv.FormatInt(v.Id, 10),
			ReconsumeTimes:  fcount,
			BornTimestamp:   ts,
			StoreTimestamp:  ts,
			QueueOffset:     v.Id,
			CommitLogOffset: v.Id,
		}

		err1 := fu(msgExt)
		if err1 != nil && err1.Code() != errors.OK.Code() {
			// 消费失败，重新消费
			log.Errorf("consumer fail: %v", msgExt)
			consumerFailMessageId = v.Id
			break
		}
		consumerMessageId = v.Id
	}

	return consumerMessageId, consumerFailMessageId
}

func checkConsumerExist(topic string, groupId string) bool {
	key := topic + "#" + groupId

	if _, ok := consumers[key]; ok {
		return true
	}
	return false
}

func queryTopicMessage(consumer *PpmMqsMessageQueueConsumer, count int) (*[]PpmMqsMessageQueue, errors.SystemErrorInfo) {
	cond := db.Cond{}
	cond[consts.TcIsDelete] = 2
	cond[consts.TcTopic] = consumer.Topic
	cond[consts.TcId] = db.Gt(consumer.MessageId)

	mqs := &[]PpmMqsMessageQueue{}

	err := mysql.SelectAllByCondWithNumAndOrder(consts.TableMessageQueue, cond, nil, 1, count, "id asc", mqs)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return nil, errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}
	return mqs, nil
}

func queryTopicGroupMessageFail(consumer *PpmMqsMessageQueueConsumer, mqs *[]PpmMqsMessageQueue) (*map[int64]PpmMqsMessageQueueConsumerFail, errors.SystemErrorInfo) {

	ids := make([]int64, 0, len(*mqs))

	for _, v := range *mqs {
		ids = append(ids, v.Id)
	}

	cond := db.Cond{}
	cond[consts.TcIsDelete] = 2
	cond[consts.TcTopic] = consumer.Topic
	cond[consts.TcGroupName] = consumer.GroupName
	cond[consts.TcMessageId] = db.In(ids)

	mqcfs := &[]PpmMqsMessageQueueConsumerFail{}

	err := mysql.SelectAllByCond(consts.TableMessageQueueConsumerFail, cond, mqcfs)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return nil, errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}

	mqcfm := &map[int64]PpmMqsMessageQueueConsumerFail{}
	if mqcfs == nil || len(*mqcfs) < 1 {
		return mqcfm, nil
	}

	for _, v := range *mqcfs {
		(*mqcfm)[v.MessageId] = v
	}

	return mqcfm, nil
}

func getTopicConsumer(topic string, groupId string) (*PpmMqsMessageQueueConsumer, errors.SystemErrorInfo) {

	cond := db.Cond{}
	cond[consts.TcIsDelete] = 2
	cond[consts.TcTopic] = topic
	cond[consts.TcGroupName] = groupId

	consumer := &PpmMqsMessageQueueConsumer{}
	err1 := mysql.SelectOneByCond(consumer.TableName(), cond, consumer)
	if err1 != nil {
		// 查询不到consumer,需创建
		lock.Lock()
		err1 = mysql.SelectOneByCond(consumer.TableName(), cond, consumer)
		if err1 != nil {
			id := snowflake.Id()
			consumer = &PpmMqsMessageQueueConsumer{
				Id:               id,
				Topic:            topic,
				GroupName:        groupId,
				MessageId:        0,
				LastConsumerTime: time.Now(),
				Status:           1,
				CreateTime:       time.Now(),
				UpdateTime:       time.Now(),
				Version:          1,
				IsDelete:         2,
			}
			err1 := mysql.Insert(consumer)
			if err1 != nil {
				return nil, errors.BuildSystemErrorInfo(errors.DbMQCreateConsumerError, err1)
			}
		}
		lock.Unlock()
	}

	b := checkConsumerExist(topic, groupId)
	if b == true {
		return nil, errors.BuildSystemErrorInfo(errors.DbMQConsumerStartedError)
	}

	lock.Lock()
	b = checkConsumerExist(topic, groupId)
	if b == true {
		return nil, errors.BuildSystemErrorInfo(errors.DbMQConsumerStartedError)
	}
	key := topic + "#" + groupId
	consumers[key] = 1
	lock.Unlock()

	return consumer, nil
}

func GetDbMQProxy() *DbMQProxy {
	return &DbMQProxy{}
}
