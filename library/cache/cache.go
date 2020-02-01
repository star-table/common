package cache

import (
	"github.com/galaxy-book/common/core/config"
	"github.com/galaxy-book/common/core/consts"
	"github.com/galaxy-book/common/library/cache/hashmap"
	"github.com/galaxy-book/common/library/cache/redis"
	"sync"
)

type Cache interface {
	Set(key string, value string) error
	SetEx(key string, value string, second int64) error
	Get(key string) (string, error)
	Del(keys ...interface{}) (int64, error)
	Exist(key string) (bool, error)
	Expire(key string, expire int64) (bool, error)
	Incrby(key string, v int64) (int64, error)
	MGet(keys ...interface{}) ([]string, error)
	MSet(kvs map[string]string) error
	HGet(key, field string) (string, error)
	HSet(key, field, value string) error
	HDel(key string, fields ...interface{}) (int64, error)
	HExists(key, field string) (bool, error)
	HMGet(key string, fields ...interface{}) (map[string]*string, error)
	HMSet(key string, fieldValue map[string]string) error
	HINCRBY(key string, field string, increment int64) (int64, error)

	TryGetDistributedLock(key string, v string) (bool, error)
	ReleaseDistributedLock(key string, v string) (bool, error)
}

var (
	insideCache Cache = &hashmap.CacheMap{
		Cache: sync.Map{},
	}
	redisCache Cache = redis.GetProxy()
)

func getCache() Cache {
	if consts.CacheModeRedis == config.GetApplication().CacheMode {
		return redisCache
	}

	return insideCache
}

func Get(key string) (string, error) {
	return getCache().Get(key)
}

func Set(key string, value string) error {
	return getCache().Set(key, value)
}

func MSet(kvs map[string]string) error{
	return getCache().MSet(kvs)
}

func SetEx(key string, value string, second int64) error {
	return getCache().SetEx(key, value, second)
}

func Del(keys ...interface{}) (int64, error) {
	return getCache().Del(keys...)
}

func Exist(key string) (bool, error) {
	return getCache().Exist(key)
}

func Expire(key string, expire int64) (bool, error) {
	return getCache().Expire(key, expire)
}

func Incrby(key string, v int64) (int64, error) {
	return getCache().Incrby(key, v)
}

func TryGetDistributedLock(key string, v string) (bool, error) {
	return getCache().TryGetDistributedLock(key, v)
}

func ReleaseDistributedLock(key string, v string) (bool, error) {
	return getCache().ReleaseDistributedLock(key, v)
}

func MGet(keys ...interface{}) ([]string, error) {
	return getCache().MGet(keys...)
}

func HGet(key, field string) (string, error) {
	return getCache().HGet(key, field)
}

func HSet(key, field, value string) error {
	return getCache().HSet(key, field, value)
}

func HDel(key string, fields ...interface{}) (int64, error) {
	return getCache().HDel(key, fields...)
}

func HExists(key, field string) (bool, error) {
	return getCache().HExists(key, field)
}

func HMGet(key string, fields ...interface{}) (map[string]*string, error) {
	return getCache().HMGet(key, fields...)
}

func HMSet(key string, fieldValue map[string]string) error {
	return getCache().HMSet(key, fieldValue)
}

func HINCRBY(key string, field string, increment int64) (int64, error) {
	return getCache().HINCRBY(key, field, increment)
}
