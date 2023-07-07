package maps

import (
	"testing"

	"github.com/star-table/common/core/util/json"
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
