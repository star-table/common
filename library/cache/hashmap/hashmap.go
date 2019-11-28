package hashmap

import (
	"gitea.bjx.cloud/allstar/common/core/lock"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"gitea.bjx.cloud/allstar/common/core/util/slice"
	"strconv"
	"sync"
)

//var cacheMap = &CacheMap{
//	Cache: sync.Map{},
//}

type CacheMap struct {
	Cache sync.Map
}

type HashType map[string]string

func (c *CacheMap) Set(key string, value string) error {
	c.Cache.Store(key, value)
	return nil
}

func (c *CacheMap) SetEx(key string, value string, milli int64) error {
	c.Cache.Store(key, value)
	return nil
}

func (c *CacheMap) Get(key string) (string, error) {
	if v, ok := c.Cache.Load(key); ok {
		return v.(string), nil
	}
	return "", nil
}

func (c *CacheMap) Del(key string) (int64, error) {
	c.Cache.Delete(key)
	return 1, nil
}

func (c *CacheMap) Exist(key string) (bool, error) {
	if _, ok := c.Cache.Load(key); ok {
		return true, nil
	}
	return false, nil
}

func (c *CacheMap) Expire(key string, expire int64) (bool, error) {
	return true, nil
}

func (c *CacheMap) Incrby(key string, v int64) (int64, error) {
	lock.Lock(key)
	defer lock.Unlock(key)
	exist, err := c.Exist(key)
	if err != nil {
		return 0, err
	}
	if !exist {
		c.Set(key, strconv.FormatInt(v, 10))
		return v, nil
	}
	value, _ := c.Cache.Load(key)
	count, err := strconv.ParseInt(value.(string), 10, 64)
	if err != nil {
		return 0, err
	}
	count += v
	c.Set(key, strconv.FormatInt(count, 10))
	return count, nil
}

func (c *CacheMap) MGet(keys ...interface{}) ([]string, error) {
	resultList := make([]string, 0)
	for _, key := range keys {
		if v, ok := c.Cache.Load(key); ok {
			resultList = append(resultList, v.(string))
		}
	}
	return resultList, nil
}

func (c *CacheMap) TryGetDistributedLock(key string, v string) (bool, error) {
	lock.Lock(key)
	return true, nil
}

func (c *CacheMap) ReleaseDistributedLock(key string, v string) (bool, error) {
	lock.Unlock(key)
	return true, nil
}

func (c *CacheMap) HGet(key, field string) (string, error) {
	mid, ok := c.Cache.Load(key)
	if !ok {
		return "", nil
	}
	mapRes := &HashType{}
	err := json.FromJson(mid.(string), mapRes)
	if err != nil {
		return "", err
	}
	for k, v := range *mapRes {
		if k == field {
			return v, nil
		}
	}
	return "", nil
}

func (c *CacheMap) HSet(key, field, value string) error {
	mid, ok := c.Cache.Load(key)
	if !ok {
		cacheValue := HashType{
			field: value,
		}
		c.Cache.Store(key, json.ToJsonIgnoreError(cacheValue))
	} else {
		mapRes := &HashType{}
		err := json.FromJson(mid.(string), mapRes)
		if err != nil {
			return err
		}
		(*mapRes)[field] = value
		c.Cache.Store(key, json.ToJsonIgnoreError(mapRes))
	}

	return nil
}

func (c *CacheMap) HExists(key, field string) (bool, error) {
	mid, ok := c.Cache.Load(key)
	if !ok {
		return false, nil
	}
	mapRes := &HashType{}
	err := json.FromJson(mid.(string), mapRes)
	if err != nil {
		return false, err
	}
	for k, _ := range *mapRes {
		if k == field {
			return true, nil
		}
	}
	return false, nil
}

func (c *CacheMap) HDel(key string, fields ...interface{}) (int64, error) {
	mid, ok := c.Cache.Load(key)
	var i int64
	if !ok {
		return 0, nil
	}
	mapRes := &HashType{}
	err := json.FromJson(mid.(string), mapRes)
	if err != nil {
		return 0, err
	}
	for k, _ := range *mapRes {
		if ok, _ = slice.Contain(fields, k); ok {
			i++
			delete(*mapRes, k)
		}
	}

	c.Cache.Store(key, json.ToJsonIgnoreError(mapRes))
	return i, nil
}

func (c *CacheMap) HMGet(key string, fields ...interface{}) (map[string]*string, error) {
	mid, ok := c.Cache.Load(key)
	result := map[string]*string{}
	if !ok {
		return result, nil
	}
	mapRes := &HashType{}
	err := json.FromJson(mid.(string), mapRes)
	if err != nil {
		return result, err
	}
	for _, v := range fields {
		result[v.(string)] = nil
	}
	for k, v := range *mapRes {
		if ok, _ = slice.Contain(fields, k); ok {
			result[k] = &v
		}
	}
	return result, nil
}
