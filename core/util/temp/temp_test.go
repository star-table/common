package temp

import (
	"fmt"
	"testing"
)

type TestStruct struct {
	Name string
	Age  int64
}

func TestRender(t *testing.T) {
	ts := TestStruct{
		Name: "Nico",
		Age:  1,
	}

	str := "这是一个测试 {{.}}"
	fmt.Println(Render(str, 1))

	str = "这是第二个模板测试 {{.Name}}{{.Age}}"
	fmt.Println(Render(str, ts))

	tmap := map[string]string{
		"Name": "1111",
	}
	str = "这是第三个模板测试 {{.Name}}"
	fmt.Println(Render(str, tmap))
}
