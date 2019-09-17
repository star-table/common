package dbmq

import "time"

type PpmMqsMessageQueue struct {
	Id					int64		`db:"id,omitempty" json:"Id"`
	Topic					string		`db:"topic,omitempty" json:"Topic"`
	MessageKey					string		`db:"message_key,omitempty" json:"MessageKey"`
	Message					string		`db:"message,omitempty" json:"Message"`
	Status					int		`db:"status,omitempty" json:"Status"`
	Creator					int64		`db:"creator,omitempty" json:"Creator"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"CreateTime"`
	Updator					int64		`db:"updator,omitempty" json:"Updator"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"UpdateTime"`
	Version					int		`db:"version,omitempty" json:"Version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"IsDelete"`
}

func (*PpmMqsMessageQueue) TableName() string {
	return "ppm_mqs_message_queue"
}

type PpmMqsMessageQueueConsumerFail struct {
	Id					int64		`db:"id,omitempty" json:"Id"`
	Topic					string		`db:"topic,omitempty" json:"Topic"`
	GroupName					string		`db:"group_name,omitempty" json:"GroupName"`
	MessageId					int64		`db:"message_id,omitempty" json:"MessageId"`
	FailCount					int		`db:"fail_count,omitempty" json:"FailCount"`
	FailTime					time.Time		`db:"fail_time,omitempty" json:"FailTime"`
	Status					int		`db:"status,omitempty" json:"Status"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"CreateTime"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"UpdateTime"`
	Version					int		`db:"version,omitempty" json:"Version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"IsDelete"`
}

func (*PpmMqsMessageQueueConsumerFail) TableName() string {
	return "ppm_mqs_message_queue_consumer_fail"
}

type PpmMqsMessageQueueConsumer struct {
	Id					int64		`db:"id,omitempty" json:"Id"`
	Topic					string		`db:"topic,omitempty" json:"Topic"`
	GroupName					string		`db:"group_name,omitempty" json:"GroupName"`
	MessageId					int64		`db:"message_id,omitempty" json:"MessageId"`
	LastConsumerTime					time.Time		`db:"last_consumer_time,omitempty" json:"LastConsumerTime"`
	Status					int		`db:"status,omitempty" json:"Status"`
	CreateTime					time.Time		`db:"create_time,omitempty" json:"CreateTime"`
	UpdateTime					time.Time		`db:"update_time,omitempty" json:"UpdateTime"`
	Version					int		`db:"version,omitempty" json:"Version"`
	IsDelete					int		`db:"is_delete,omitempty" json:"IsDelete"`
}

func (*PpmMqsMessageQueueConsumer) TableName() string {
	return "ppm_mqs_message_queue_consumer"
}


