package redis

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestAvailability(t *testing.T) {
	t.Logf("start load config")
	config.GetConfig().Redis = &config.RedisConfig{
		Host:       "172.19.132.101",
		Port:       26379,
		IsSentinel: true,
		MasterName: "mymaster",
		Database:   0,
	}
	rp := GetProxy()
	err := rp.Set("abc", "abababba")
	t.Log(err)
	v, err := rp.Get("abc")
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

func Test_Pool(t *testing.T) {
	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\config", "application")
	for i := 0; i < 1000; i++ {
		go getConn()
	}
	fmt.Print("阻塞")

	time.Sleep(time.Duration(10) * time.Second)
}

func getConn() {
	conn := GetRedisConn()
	time.Sleep(time.Duration(1) * time.Second)

	conn.Close()
}

func Test_MGet(t *testing.T) {
	config.LoadUnitTestConfig()
	rp := GetProxy()

	rp.Set("3", "abc")
	rp.Set("4", "efg")
	v, e := rp.MGet("3", "4")
	t.Log(e)
	t.Log(v)
}

func TestProxy_HGet(t *testing.T) {
	config.LoadUnitTestConfig()
	rp := GetProxy()

	key := "abc"
	t.Log(rp.HSet(key, "a", "a"))
	t.Log(rp.HSet(key, "b", "b"))
	t.Log(rp.HDel(key, "a"))
	t.Log(rp.HGet(key, "a"))
	res, err := rp.HMGet(key, "a", "b", "c", "a")
	t.Log(json.ToJsonIgnoreError(res), err)
	assert.Equal(t, err, nil)
	t.Log(rp.HMSet(key, map[string]string{
		"h": "h",
		"i": "i",
	}))
	res1, err := rp.HMGet(key, "a", "b", "c", "h", "i")
	t.Log(json.ToJsonIgnoreError(res1), err)
	assert.Equal(t, err, nil)

}

func TestProxy_Del(t *testing.T) {
	config.LoadUnitTestConfig()
	rp := GetProxy()

	t.Log(rp.Set("a", "a"))
	t.Log(rp.Set("b", "b"))
	t.Log(rp.Set("c", "c"))
	t.Log(rp.MGet("a", "b", "c"))
	t.Log(rp.Del("a"))
	t.Log(rp.Del("b", "c"))
	t.Log(rp.MGet("a", "b", "c"))
}

func TestProxy_Exist(t *testing.T) {
	config.LoadUnitTestConfig()
	rp := GetProxy()

	//key := "abc"
	//b := int64(1)
	//t.Log(rp.HSet(key, strconv.FormatInt(b, 10), "bbbbb"))
	//t.Log(rp.HGet(key, strconv.FormatInt(b, 10)))
	//t.Log(rp.HDel(key, 1))
	//t.Log(rp.HGet(key, strconv.FormatInt(b, 10)))
	t.Log(rp.HSet("polaris:rolesvc:org_10101:user_10201:user_role_list_hash:1", strconv.FormatInt(0, 10), "11"))
	t.Log(rp.HGet("polaris:rolesvc:org_10101:user_10201:user_role_list_hash:1", strconv.FormatInt(0, 10)))
	t.Log(rp.HDel("polaris:rolesvc:org_10101:user_10201:user_role_list_hash:1", strconv.FormatInt(0, 10)))
	t.Log(rp.HGet("polaris:rolesvc:org_10101:user_10201:user_role_list_hash:1", strconv.FormatInt(0, 10)))
}

func TestProxy_HINCRBY(t *testing.T) {
	config.LoadUnitTestConfig()
	rp := GetProxy()
	t.Log(rp.HINCRBY("aaa", "bbb", 1))
}
