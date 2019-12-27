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

type HashType map[string]string

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
	resultList := make([]string, 0)
	for _, v := range list {
		if bytes, ok := v.([]byte); ok {
			resultList = append(resultList, string(bytes))
		}
	}
	return resultList, nil
}

func (rp *Proxy) Del(keys ...interface{}) (int64, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return 0, e
	}
	rs, err := conn.Do("Del", keys...)
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

func (rp *Proxy) HGet(key, field string) (string, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return "", e
	}
	rs, err := conn.Do("HGET", key, field)
	if err != nil {
		return "", err
	}
	if rs == nil {
		return "", err
	}
	return string(rs.([]byte)), nil
}

func (rp *Proxy) HSet(key, field, value string) error {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return e
	}
	_, err := conn.Do("HSET", key, field, value)
	return err
}

func (rp *Proxy) HDel(key string, fields ...interface{}) (int64, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return 0, e
	}
	args := []interface{}{}
	args = append(append(args, key), fields...)
	rs, err := conn.Do("HDEL", args...)
	if err != nil {
		return 0, err
	}
	if rs == nil {
		return 0, err
	}
	return rs.(int64), nil
}

func (rp *Proxy) HExists(key, field string) (bool, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return false, e
	}
	rs, err := conn.Do("HEXISTS", key, field)
	if err != nil {
		return false, err
	}
	if rs == nil {
		return false, err
	}
	return rs.(int64) == 1, err
}

func (rp *Proxy) HMGet(key string, fields ...interface{}) (map[string]*string, error) {
	result := map[string]*string{}
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return result, e
	}
	fields = append([]interface{}{key}, fields...)
	rs, err := conn.Do("HMGET", fields...)
	if err != nil {
		return result, err
	}
	if rs == nil {
		return result, nil
	}
	keys := make([]string, len(fields)-1)
	for i, k := range fields {
		if i == 0 {
			continue
		}
		keys[i-1] = k.(string)
	}

	if datas, ok := rs.([]interface{}); ok {
		for i, data := range datas {
			if data == nil {
				result[keys[i]] = nil
			} else {
				dataStr := string(data.([]byte))
				result[keys[i]] = &dataStr
			}
		}
	}
	return result, nil
}

func (rp *Proxy) HMSet(key string, fieldValue map[string]string) error {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return e
	}
	args := []interface{}{}
	args = append(args, key)
	for k, v := range fieldValue {
		args = append(append(args, k), v)
	}
	_, err := conn.Do("HMSET", args...)
	return err
}

func (rp *Proxy) HINCRBY(key string, field string, increment int64) (int64, error) {
	conn, e := Connect()
	defer Close(conn)
	if e != nil {
		return 0, e
	}

	res, err := conn.Do("HINCRBY", key, field, increment)
	if err != nil {
		return 0, e
	}
	return res.(int64), nil
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
