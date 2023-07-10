package config

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/star-table/common/core/util/file"
	json2 "github.com/star-table/common/core/util/json"
)

/**
加载本地配置
*/
func TestLoadUnitTestConfig2(t *testing.T) {

	convey.Convey("Test LoadLocalConfig", t, func() {
		convey.Convey("test", func() {
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

/**
添加GetOfficeUrl的单元测试
*/
//func TestGetOfficeUrl(t *testing.T) {
//	convey.Convey("Test GetOffice", t, func() {
//		// LoadUnitTestConfigWithEnv("local")
//		convey.So(GetOfficeUrl(), convey.ShouldNotBeNil)
//	})
//}
//
//func TestLoadConfig(t *testing.T) {
//
//	// //
//	// LoadNacosConfig("172.19.132.101", 8848, "polaris-common", "DEFAULT_GROUP", "application.common.yaml")
//	// LoadNacosConfig("127.0.0.1", 8848, "public", "POLARIS_COMMON", "application.common.yaml")
//	// LoadNacosConfig("127.0.0.1", 8848, "public", "POLARIS_ORGSVC", "application.common.yaml")
//	// LoadNacosConfig("127.0.0.1", 8848, "public", "POLARIS_ORGSVC", "application.common.yaml")
//
//	err := LoadNacosConfig("172.19.132.101", 8848, "polaris-common", "DEFAULT_GROUP", "application.common.yaml", "nacos", "oY4hMsdsqaB06kWK")
//	if err != nil {
//		return
//	}
//
//	fmt.Println("---------------application.common.yaml----------------")
//	fmt.Println(json2.ToJsonIgnoreError(GetConfig()))
//
//	err = LoadNacosConfig("172.19.132.101", 8848, "4c635db8-0f29-4b8a-8271-fd4881baec0b", "DEFAULT_GROUP", "application.common.gray.yaml", "nacos", "oY4hMsdsqaB06kWK")
//	if err != nil {
//		return
//	}
//
//	fmt.Println("---------------application.common.gray.yaml----------------")
//	fmt.Println(json2.ToJsonIgnoreError(GetConfig()))
//
//	err = LoadNacosConfig("172.19.132.101", 8848, "polaris-common", "DEFAULT_GROUP", "application.resource.yaml", "nacos", "oY4hMsdsqaB06kWK")
//	if err != nil {
//		return
//	}
//
//	fmt.Println("---------------application.resource.yaml----------------")
//	fmt.Println(json2.ToJsonIgnoreError(GetConfig()))
//
//	err = LoadNacosConfig("172.19.132.101", 8848, "4c635db8-0f29-4b8a-8271-fd4881baec0b", "DEFAULT_GROUP", "application.resource.yaml", "nacos", "oY4hMsdsqaB06kWK")
//	if err != nil {
//		return
//	}
//
//	fmt.Println("---------------application.resource.yaml----------------")
//	fmt.Println(json2.ToJsonIgnoreError(GetConfig()))
//}
//
//func TestLoadConfigFromNacos(t *testing.T) {
//
//	LoadNacosConfigAutoExtends("172.19.132.101", 8848, "4c635db8-0f29-4b8a-8271-fd4881baec0b", "resource", "gray", "nacos", "oY4hMsdsqaB06kWK")
//
//	fmt.Println(json2.ToJsonIgnoreError(GetConfig()))
//}
