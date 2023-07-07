package tests

import (
	"fmt"
	"github.com/polaris-team/dingtalk-sdk-golang/json"
	"github.com/star-table/common/core/config"
)

func StartUp(f func()) func() {
	return func() {
		//加载测试配置
		config.LoadUnitTestConfig()

		testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())

		fmt.Println("unitteset Mysql配置json:", testMysqlJson)

		testRedisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("unitteset redis配置json:", testRedisJson)

		//丢进来的方法立刻执行
		f()
	}
}
