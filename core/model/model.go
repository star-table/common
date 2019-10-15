package model

import "net/http"

type MqMessageExt struct {
	MqMessage

	/**
	消息id
	*/
	MessageID string
	/**
	仅rocketmq
	*/
	QueueId int
	/**
	重新消费次数
	*/
	ReconsumeTimes int
	/**
	仅rocketmq
	*/
	StoreSize int
	/**
	消息创建时间戳
	*/
	BornTimestamp int64
	/**
	持久化时间戳
	*/
	StoreTimestamp int64
	/**
	仅rocketmq
	*/
	QueueOffset int64
	/**
	仅rocketmq
	*/
	CommitLogOffset int64
	/**
	仅rocketmq
	*/
	PreparedTransactionOffset int64

	/**
	发送结果，1成功，2失败
	*/
	SendStatus int
}

type MqMessage struct {
	// 消息结果 1，待发送，2发送中，3发送成功，4发送失败
	SendResult int

	/**
	topic名
	*/
	Topic string

	/**
	rocketmq 特有
	*/
	Tags string

	/**
	消息id
	*/
	Keys string
	/**
	消息体
	*/
	Body string
	/**
	推迟处理次数
	*/
	DelayTimeLevel int

	/**

	 */
	Property map[string]string

	//消费重试次数
	ReconsumeTimes *int
	//推送重试次数
	RePushTimes *int


	//Kafka特有
	Partition int32
	Offset int64
}

type HttpContext struct {
	Request   *http.Request
	Url       string
	Method    string
	StartTime int64
	EndTime   int64
	TraceId   string
	Ip        string
	Status    int
}