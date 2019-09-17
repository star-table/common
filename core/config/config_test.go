package config

import (
	"gitea.bjx.cloud/allstar/common/core/util/file"
	json2 "gitea.bjx.cloud/allstar/common/core/util/json"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

/**
加载本地配置
*/
func TestLoadLocalConfig(t *testing.T) {

	convey.Convey("Test LoadLocalConfig", t, func() {
		convey.Convey("test", func() {
			LoadLocalConfig()
		})
	})

	bs, _ := json2.ToJson(conf)
	t.Log(string(bs))

	t.Log(file.GetCurrentPath())
}

/**
加载单元测试的配置
*/
func TestLoadUnitTestConfig(t *testing.T) {

	convey.Convey("Test LoadUnitTestConfig", t, func() {
		convey.Convey("test", func() {
			LoadUnitTestConfig()
		})
	})

	bs, _ := json2.ToJson(conf)
	t.Log(string(bs))

	t.Log(file.GetCurrentPath())
}

/**
加载单元测试的配置
*/
func TestGetEnv(t *testing.T) {
	convey.Convey("Test LoadUnitTestConfig", t, func() {

		convey.So(GetEnv(), convey.ShouldNotBeNil)
	})
}

/**
添加kafka的单元测试
*/
func TestGetKafkaConfig(t *testing.T) {

	convey.Convey("Test LoadUnitTestConfig", t, func() {

		convey.So(GetKafkaConfig(), convey.ShouldNotBeNil)
	})

}

/**
添加redis的单元测试
*/
func TestGetRedisConfig(t *testing.T) {

	convey.Convey("Test GetRedisConfig", t, func() {

		convey.So(GetRedisConfig(), convey.ShouldNotBeNil)
	})

}

/**
添加config的单元测试
*/
func TestGetConfig(t *testing.T) {

	convey.Convey("Test GetConfig", t, func() {

		convey.So(GetConfig(), convey.ShouldNotBeNil)
	})
}

/**
添加GetMailConfig的单元测试
*/
func TestGetMailConfig(t *testing.T) {

	convey.Convey("Test GetMailConfig", t, func() {

		convey.So(GetMailConfig(), convey.ShouldNotBeNil)
	})
}

/**
添加LoadConfig的单元测试
*/
func TestLoadConfig(t *testing.T) {

	convey.Convey("Test LoadConfig", t, func() {

		convey.So(LoadConfig(), convey.ShouldNotBeNil)
	})
}

/**
添加GetServerConfig的单元测试
*/
func TestGetServerConfig(t *testing.T) {

	convey.Convey("Test GetServerConfig", t, func() {

		convey.So(GetServerConfig(), convey.ShouldNotBeNil)
	})
}

/**
添加GetServerConfig的单元测试
*/
func TestGetDingTalkSdkConfig(t *testing.T) {

	convey.Convey("Test GetDingTalkSdkConfig", t, func() {

		convey.So(GetDingTalkSdkConfig(), convey.ShouldNotBeNil)
	})
}

/**
添加GetMQ的单元测试
*/
func TestGetMQ(t *testing.T) {

	convey.Convey("Test GetMQ", t, func() {

		convey.So(GetMQ(), convey.ShouldNotBeNil)
	})
}

/**
添加GetServerConfig的单元测试
*/
func TestGetElasticSearch(t *testing.T) {

	convey.Convey("Test GetElasticSearch", t, func() {

		convey.So(GetElasticSearch(), convey.ShouldNotBeNil)
	})
}
