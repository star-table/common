package mysql

const (
	OssConfigMissingSentence = "oss configuration is missing!"
	DBOpenErrorSentence      = "db.Open(): %q\n"
	TxOpenErrorSentence      = "tx.Open(): %q\n"
)

const(
	TableMessageQueue             = "ppm_mqs_message_queue"
	TableMessageQueueConsumer     = "ppm_mqs_message_queue_consumer"
	TableMessageQueueConsumerFail = "ppm_mqs_message_queue_consumer_fail"
)

const(
	TcId                       = "id"
	TcIsDelete                 = "is_delete"
	TcTopic                    = "topic"
	TcGroupName                = "group_name"
	TcMessageId                = "message_id"
)