package redis

import (
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"math/rand"
	"testing"
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
