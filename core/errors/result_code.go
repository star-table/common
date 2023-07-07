package errors

var (
	//成功
	OK = AddResultCodeInfo(0, "OK", "ResultCode.OK")

	//token错误
	RequestError = AddResultCodeInfo(400, "请求错误", "ResultCode.RequestError")

	//认证错误
	Unauthorized = AddResultCodeInfo(401, "认证错误", "ResultCode.Unauthorized")

	//禁止访问
	ForbiddenAccess = AddResultCodeInfo(403, "禁止访问", "ResultCode.ForbiddenAccess")

	//请求地址不存在
	PathNotFound = AddResultCodeInfo(404, "请求地址不存在", "ResultCode.PathNotFound")

	//不支持该方法
	MethodNotAllowed = AddResultCodeInfo(405, "不支持该方法", "ResultCode.MethodNotAllowed")

	//Token过期
	TokenExpires = AddResultCodeInfo(450, "登录失效", "ResultCode.TokenExpires")

	//请求参数错误
	ServerError = AddResultCodeInfo(500, "服务器错误", "ResultCode.ServerError")

	//过载保护,服务暂不可用
	ServiceUnavailable = AddResultCodeInfo(503, "过载保护,服务暂不可用", "ResultCode.ServiceUnavailable")

	//服务调用超时
	Deadline = AddResultCodeInfoWithSentry(504, "服务调用超时", "ResultCode.Deadline")

	//超出限制
	LimitExceed = AddResultCodeInfo(509, "超出限制", "ResultCode.LimitExceed")

	//参数错误
	ParamError = AddResultCodeInfo(600, "参数错误", "ResultCode.ParamError")

	//文件过大
	FileTooLarge = AddResultCodeInfo(610, "文件过大", "ResultCode.FileTooLarge")

	//文件类型错误
	FileTypeError = AddResultCodeInfo(611, "文件类型错误", "ResultCode.FileTypeError")

	//文件或目录不存在
	FileNotExist = AddResultCodeInfo(612, "文件或目录不存在", "ResultCode.FileNotExist")

	//文件路径为空
	FilePathIsNull = AddResultCodeInfo(613, "文件路径为空", "ResultCode.FilePathIsNull")

	//读取文件失败
	FileReadFail = AddResultCodeInfo(614, "读取文件失败", "ResultCode.FileReadFail")

	//错误未定义
	ErrorUndefined = AddResultCodeInfo(996, "错误未定义", "ResultCode.ErrorUndefined")

	//业务失败
	BusinessFail = AddResultCodeInfo(997, "业务失败", "ResultCode.BusinessFail")

	//系统异常
	SystemError = AddResultCodeInfo(998, "系统异常", "ResultCode.SystemError")

	//未知错误
	UnknownError = AddResultCodeInfo(999, "未知错误", "ResultCode.UnknownError")

	//Panic错误
	PanicError = AddResultCodeInfoWithSentry(1000, "Panic错误", "ResultCode.PanicError")

	//数据库错误
	DatabaseError = AddResultCodeInfoWithSentry(1001, "Database错误", "ResultCode.DatabaseError")

	RocketMQProduceInitError   = AddResultCodeInfo(10001, "RocketMQ Produce 初始化异常", "ResultCode.RocketMQProduceInitError")
	RocketMQSendMsgError       = AddResultCodeInfo(10002, "RocketMQ SendMsg 失败", "ResultCode.RocketMQSendMsgError")
	RocketMQConsumerInitError  = AddResultCodeInfo(10003, "RocketMQ Consumer 初始化异常", "ResultCode.RocketMQConsumerInitError")
	RocketMQConsumerStartError = AddResultCodeInfo(10004, "RocketMQ Consumer 启动异常", "ResultCode.RocketMQConsumerStartError")
	RocketMQConsumerStopError  = AddResultCodeInfo(10005, "RocketMQ Consumer 停止异常", "ResultCode.RocketMQConsumerStopError")

	KafkaMqSendMsgError           = AddResultCodeInfo(10101, "Kafka发送消息失败", "ResultCode.KafkaMqSendMsgError")
	KafkaMqSendMsgCantBeNullError = AddResultCodeInfo(10102, "Kafka发送的消息不能为空", "ResultCode.KafkaMqSendMsgCantBeNullError")
	KafkaMqConsumeMsgError        = AddResultCodeInfo(10103, "Kafka消费消息失败", "ResultCode.KafkaMqConsumeMsgError")
	KafkaMqConsumeStartError      = AddResultCodeInfo(10104, "Kafka消费启动失败", "ResultCode.KafkaMqConsumeStartError")

	TryDistributedLockError = AddResultCodeInfo(10201, "获取分布式锁异常", "ResultCode.TryDistributedLockError")
	GetDistributedLockError = AddResultCodeInfo(10202, "获取分布式锁失败", "ResultCode.GetDistributedLockError")
	MysqlOperateError       = AddResultCodeInfo(10203, "db操作出现异常", "ResultCode.MysqlOperateError")
	RedisOperateError       = AddResultCodeInfo(10204, "redis操作出现异常", "ResultCode.RedisOperateError")

	DbMQSendMsgError         = AddResultCodeInfo(10301, "Db 保存 message queue 失败", "ResultCode.DbMQSendMsgError")
	DbMQCreateConsumerError  = AddResultCodeInfo(10302, "Db 创建 message queue consumer 失败", "ResultCode.DbMQCreateConsumerError")
	DbMQConsumerStartedError = AddResultCodeInfo(10303, "Db 创建 message queue consumer 已启动", "ResultCode.DbMQConsumerStartedError")

	MQTTNoClientError = AddResultCodeInfo(10401, "Mqtt 无可用Client", "ResultCode.MQTTNoClientError")
	MQTTNoConfigError = AddResultCodeInfo(10402, "Mqtt 无配置", "ResultCode.MQTTNoConfigError")
)
