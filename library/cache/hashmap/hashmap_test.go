package hashmap

import (
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCacheMap_MGet(t *testing.T) {

	c := CacheMap{}

	c.Set("1", "ab")
	c.Set("2", "ef")

	res, err := c.MGet("1", "2")
	t.Log(err)
	t.Log(res)

}

func TestCacheMap_HGet(t *testing.T) {
	c := CacheMap{}
	key := map[string]string{
		"a": "a",
		"b": "b",
	}
	c.Set("abc", json.ToJsonIgnoreError(key))
	t.Log(c.HGet("abc", "a"))
}

func TestCacheMap_HSet(t *testing.T) {
	c := CacheMap{}
	key := "aaa"
	t.Log(c.HSet(key, "a", "a"))
	t.Log(c.HSet(key, "b", "b"))
	t.Log(c.HGet(key, "a"))
	t.Log(c.HGet(key, "c"))
	t.Log(c.HExists(key, "a"))
	t.Log(c.HExists(key, "c"))
	t.Log(c.HDel(key, "a", "c"))
	t.Log(c.HGet(key, "b"))
	t.Log(c.HGet(key, "a"))
	res, err := c.HMGet(key, "a", "b", "c")
	t.Log(json.ToJsonIgnoreError(res), err)
	assert.Equal(t, err, nil)
	t.Log(c.HMSet(key, map[string]string{
		"h": "h",
		"i": "i",
	}))
	res1, err := c.HMGet(key, "a", "b", "c", "h", "i")
	t.Log(json.ToJsonIgnoreError(res1), err)
	assert.Equal(t, err, nil)
}
