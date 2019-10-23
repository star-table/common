package consts

import "time"

//Token
const AppHeaderTokenName = "Token"

// linux操作系统
const LinuxGOOS = "linux"

//默认空时间
const BlankTime = "1970-01-01 00:00:00"
const BlankDate = "1970-01-01"

var BlankTimeObject, _ = time.Parse(AppTimeFormat, BlankTime)

//默认空字符串
const BlankString = ""

//是否被删除
const (
	AppIsDeleted  = 1
	AppIsNoDelete = 2
)

const (
	AppUserIsInUse     = 1
	AppUserIsNotInUser = 2
)

//是否流程初始化状态
const (
	//是
	AppIsInitStatus = 1
	//否
	AppIsNotInitStatus = 2
)

//是否可用
const (
	AppStatusEnable   = 1
	AppStatusDisabled = 2
)

//是否默认
const (
	APPIsDefault    = 1
	AppIsNotDefault = 2
)

//全局日期格式
const AppDateFormat = "2006-01-02"
const AppTimeFormat = "2006-01-02 15:04:05"
const AppSystemTimeFormat = "2006-01-02T15:04:05Z"
const AppSystemTimeFormat8 = "2006-01-02T15:04:05+08:00"

const (
	// SAAS运行模式
	AppRunmodeSaas = 1
	// 单机部署
	AppRunmodeSingle = 2
	// 私有化部署
	AppRunmodePrivate = 3
)

//初始化时的一些常量定义
const (
	InitDefaultTeamName     = "默认团队"
	InitDefaultTeamNickname = "默认团队昵称"
)

// context key
const (
	TraceIdKey     = "PM-TRACE-ID"
	HttpContextKey = "_httpContext"
)

// 默认对象id步长
const (
	DefaultObjectIdStep = 200
)

// 系统缓存模式
const (
	CacheModeRedis  = "Redis"
	CacheModeInside = "Inside"
)

// 系统消息队列模式
const (
	MQModeRocketMQ = "RocketMQ"
	MQModeDB       = "DB"
	MQModeKafka    = "Kafka"
)

// 发送消息状态
const (
	SendMQStatusSuccess = 1
	SendMQStatusFail    = 2
)

// 消息处理状态
const (
	//待处理
	MQStatusWait = 1
	//处理中
	MQStatusHandle = 2
	//处理成功
	MQStatusSuccess = 3
	//处理失败
	MQStatusFail = 4
)

const (
	OssConfigMissingSentence = "oss configuration is missing!"
	DBOpenErrorSentence      = "db.Open(): %q\n"
	TxOpenErrorSentence      = "tx.Open(): %q\n"
)

const (
	TableMessageQueue             = "ppm_mqs_message_queue"
	TableMessageQueueConsumer     = "ppm_mqs_message_queue_consumer"
	TableMessageQueueConsumerFail = "ppm_mqs_message_queue_consumer_fail"
)

const (
	TcId        = "id"
	TcIsDelete  = "is_delete"
	TcTopic     = "topic"
	TcGroupName = "group_name"
	TcMessageId = "message_id"
)
