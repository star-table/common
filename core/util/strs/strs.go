package strs

import "fmt"

func Len(str string) int {
	return len([]rune(str))
}

func ObjectToString(obj interface{}) string {
	return fmt.Sprintf("%#v", obj)
}
