module github.com/star-table/common

go 1.13

require (
	github.com/Shopify/sarama v1.29.1
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.18
	github.com/aliyun/aliyun-oss-go-sdk v2.0.1+incompatible
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/bwmarrin/snowflake v0.3.0
	github.com/bytedance/go-tagexpr v2.2.0+incompatible
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/getsentry/sentry-go v0.12.0
	github.com/go-kratos/kratos/contrib/config/nacos/v2 v2.0.0-20230706115902-bffc1a0989a6
	github.com/go-kratos/kratos/v2 v2.6.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/henrylee2cn/goutil v0.0.0-20190701081944-69b028c17642 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/jtolds/gls v4.20.0+incompatible
	github.com/magiconair/properties v1.8.5
	github.com/mozillazg/go-pinyin v0.15.0
	github.com/nacos-group/nacos-sdk-go v1.0.9
	github.com/nyaruka/phonenumbers v1.0.43 // indirect
	github.com/olivere/elastic/v7 v7.0.6
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/polaris-team/converter v0.0.8
	github.com/polaris-team/dingtalk-sdk-golang v0.0.9
	github.com/qustavo/sqlhooks/v2 v2.1.0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cast v1.4.1
	github.com/spf13/viper v1.10.0
	github.com/star-table/emitter-go-client v0.0.0-20230706083051-a677df920eca
	github.com/star-table/go-common v1.0.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/zap v1.21.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.3.3
	gorm.io/gorm v1.23.5
	upper.io/db.v3 v3.7.1+incompatible
)

replace upper.io/db.v3 v3.7.1+incompatible => github.com/star-table/db v0.3.75-0.20230707012646-28b2e2303a74
