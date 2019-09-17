package maps

import (
	"gitea.bjx.cloud/allstar/common/core/util/slice"
	"reflect"
)

type LocalMap map[interface{}]interface{}

func NewMap(key string, source interface{}) LocalMap{
	list := slice.ToSlice(source)
	localMap := LocalMap{}
	for _, obj := range list {
		immutable := reflect.ValueOf(obj)
		v := immutable.FieldByName(key).Interface()
		localMap[v] = obj
	}
	return localMap
}