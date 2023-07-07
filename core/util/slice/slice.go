package slice

import (
	"errors"
	"reflect"
)

//通过map去重slice
func SliceUniqueString(s []string) []string {
	res := make([]string, 0)
	exist := make(map[string]bool)
	for _, s2 := range s {
		if _, ok := exist[s2]; ok {
			continue
		}

		res = append(res, s2)
		exist[s2] = true
	}

	return res
}

func SliceUniqueInt64(s []int64) []int64 {
	res := make([]int64, 0)
	exist := make(map[int64]bool)
	for _, i2 := range s {
		if _, ok := exist[i2]; ok {
			continue
		}
		res = append(res, i2)
		exist[i2] = true
	}

	return res
}

func Contain(list interface{}, obj interface{}) (bool, error) {
	if list == nil {
		return false, nil
	}
	targetValue := reflect.ValueOf(list)
	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
		return false, nil
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("not in array")
}

func ToSlice(arr interface{}) []interface{} {
	if arr == nil {
		return []interface{}{}
	}
	v := reflect.ValueOf(arr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}
