package _map

import (
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"testing"
)

type TestStruct struct {
	Id   int64
	Name string
}

func TestNewCacheMap(t *testing.T) {
	list := &[]TestStruct{
		{
			Id:   1,
			Name: "hello",
		},
		{
			Id:   2,
			Name: "world",
		},
	}

	cacheMap := NewMap("Id", list)
	t.Log(json.ToJsonIgnoreError(cacheMap))
	t.Log(cacheMap)
	t.Log(cacheMap[int64(1)])
}
