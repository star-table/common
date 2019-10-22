package hashmap

import "testing"

func TestCacheMap_MGet(t *testing.T) {

	c := CacheMap{}

	c.Set("1", "ab")
	c.Set("2", "ef")

	res, err := c.MGet("1", "2")
	t.Log(err)
	t.Log(res)

}