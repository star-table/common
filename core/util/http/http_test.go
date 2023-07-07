package http

import (
	"testing"
	"time"
)

type TestStruct struct {
	Page *int
	Size *int
}

func TestConvertToQueryParams(t *testing.T) {
	ts := TestStruct{}
	data := map[string]interface{}{
		"a":  "aa",
		"b":  123434343232332,
		"bb": 1.12313213131313,
		"c":  int64(5),
		"e":  true,
		"t":  time.Now(),
		//"t1":   types.Time(time.Now()),
		"page": ts.Page,
		"size": ts.Size,
	}
	str := ConvertToQueryParams(data)
	t.Log(str)
	//assert.Equal(t, str, "?e=true&a=aa&b=123&c=5")
}
