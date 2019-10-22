package cache

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/library/cache/hashmap"
	"gitea.bjx.cloud/allstar/common/library/cache/redis"
	"sync"
)

type Cache interface {
	Set(key string, value string) error
	SetEx(key string, value string, second int64) error
	Get(key string) (string, error)
	Del(key string) (int64, error)
	Exist(key string) (bool, error)
	Expire(key string, expire int64) (bool, error)
	Incrby(key string, v int64) (int64, error)
	MGet(keys ...interface{}) ([]string, error)

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

func SetEx(key string, value string, second int64) error {
	return getCache().SetEx(key, value, second)
}

func Del(key string) (int64, error) {
	return getCache().Del(key)
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
