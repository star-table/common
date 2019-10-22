package redis

import (
	"gitea.bjx.cloud/allstar/common/core/util/times"
	"math/rand"
	"reflect"
	"sync"

	"github.com/gomodule/redigo/redis"
)

var mu sync.Mutex
var single Proxy

const (
	_LockDistributedLua = "local v;" +
		"v = redis.call('setnx',KEYS[1],ARGV[1]);" +
		"if tonumber(v) == 1 then\n" +
		"    redis.call('expire',KEYS[1],ARGV[2])\n" +
		"end\n" +
		"return v"

	_UnLockDistributedLua = "if redis.call('get',KEYS[1]) == ARGV[1]\n" +
		"then\n" +
		"    return redis.call('del',KEYS[1])\n" +
		"else\n" +
		"    return 0\n" +
		"end"

	_DistributedTimeOut = 4
	_DistributedSuccess = 1
)

var (
	_LockDistributedLuaScript   = redis.NewScript(1, _LockDistributedLua)
	_UnLockDistributedLuaScript = redis.NewScript(1, _UnLockDistributedLua)
)

type Proxy struct {
	conn redis.Conn
}

func (rp *Proxy) ZAdd(key string, score float64, value string) error {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return e
	}
	_, err := conn.Do("zadd", key, score, value)
	return err
}

func (rp *Proxy) SetEx(key string, value string, ex int64) error {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return e
	}
	_, err := conn.Do("setex", key, ex, value)
	return err
}

func (rp *Proxy) Set(key string, value string) error {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return e
	}
	_, err := conn.Do("set", key, value)
	return err
}

func (rp *Proxy) Get(key string) (string, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return "", e
	}
	rs, err := conn.Do("get", key)
	if err != nil {
		return "", err
	}
	if rs == nil {
		return "", err
	}
	return string(rs.([]byte)), nil
}

func (rp *Proxy) MGet(keys ...interface{}) ([]string, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return nil, e
	}
	rs, err := conn.Do("MGET", keys...)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, err
	}
	list := rs.([]interface{})
	resultList := make([]string, len(list))
	for i, v := range list{
		resultList[i] = string(v.([]byte))
	}
	return resultList, nil
}

func (rp *Proxy) Del(key string) (int64, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return 0, e
	}
	rs, err := conn.Do("Del", key)
	if err != nil {
		return 0, err
	}
	if rs == nil {
		return 0, err
	}
	return rs.(int64), nil
}

func (rp *Proxy) Incrby(key string, v int64) (int64, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return 0, e
	}
	rs, err := conn.Do("INCRBY", key, v)
	if err != nil {
		return 0, err
	}
	if rs == nil {
		return 0, err
	}
	return rs.(int64), err
}

func (rp *Proxy) Exist(key string) (bool, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return false, e
	}
	rs, err := conn.Do("EXISTS", key)
	if err != nil {
		return false, err
	}
	if rs == nil {
		return false, err
	}
	return rs.(int64) == 1, err
}

func (rp *Proxy) Expire(key string, expire int64) (bool, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return false, e
	}
	rs, err := conn.Do("EXPIRE", key, expire)
	if err != nil {
		return false, err
	}
	if rs == nil {
		return false, err
	}
	return rs.(int64) == 1, err
}

func (rp *Proxy) TryGetDistributedLock(key string, v string) (bool, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return false, e
	}

	end := times.GetNowMillisecond() + _DistributedTimeOut*1000
	for times.GetNowMillisecond() <= end {
		rs, err := _LockDistributedLuaScript.Do(conn, key, v, _DistributedTimeOut)
		if err != nil {
			return false, err
		}
		if rs.(int64) == _DistributedSuccess {
			return true, nil
		}
		times.SleepMillisecond(80 + int64(rand.Int31n(30)))
	}

	return false, nil
}

func (rp *Proxy) ReleaseDistributedLock(key string, v string) (bool, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return false, e
	}

	rs, err := _UnLockDistributedLuaScript.Do(conn, key, v)
	if err != nil {
		return false, err
	}
	if rs.(int64) == _DistributedSuccess {
		return true, nil
	}

	return false, nil
}

func Connect() (redis.Conn, error) {
	return GetRedisConn(), nil
}

func Close(conn redis.Conn) error {
	if conn != nil && conn.Err() == nil {
		conn.Close()
		return conn.Close()
	}
	return nil
}

func (rp Proxy) IsEmpty() bool {
	return reflect.DeepEqual(rp, Proxy{})
}

//GetProxy get redis oper proxy
func GetProxy() *Proxy {
	return &Proxy{}
}
