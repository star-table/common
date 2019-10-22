package redis

import (
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"math/rand"
	"testing"
	"time"
)

func TestAvailability(t *testing.T) {
	t.Logf("start load config")
	config.LoadUnitTestConfig()

	rp := GetProxy()

	v, err := rp.Get("abcdaf")
	t.Log(err)
	t.Log(v)

}

func TestProxy_TryGetDistributedLock(t *testing.T) {

	config.LoadUnitTestConfig()
	rp := GetProxy()

	v, err := rp.TryGetDistributedLock("aaaaa", "1333")
	t.Log(err)
	t.Log(v)

	fmt.Println("第二次获取")
	v, err = rp.TryGetDistributedLock("aaaaa", "1333")
	t.Log(err)
	t.Log(v)
	fmt.Println("第二次获取锁等待完毕")

	v, err = rp.ReleaseDistributedLock("aaaaa", "1333")
	t.Log(err)
	t.Log(v)

	for i := 1; i < 100; i++ {
		fmt.Println(rand.Int31n(30))
	}
}

func Test_Pool(t *testing.T){
	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\config", "application")
	for i := 0; i < 1000; i ++{
		go getConn()
	}
	fmt.Print("阻塞")

	time.Sleep(time.Duration(10) * time.Second)
}

func getConn(){
	conn := GetRedisConn()
	time.Sleep(time.Duration(1) * time.Second)

	conn.Close()
}

func Test_MGet(t *testing.T){
	config.LoadUnitTestConfig()
	rp := GetProxy()

	rp.Set("3", "abc")
	rp.Set("4", "efg")
	v, e := rp.MGet("3", "4")
	t.Log(e)
	t.Log(v)
}