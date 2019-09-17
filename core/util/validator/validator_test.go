package validator

import (
	"fmt"
	"testing"
)

func TestValidate(ts *testing.T) {
	type TT struct {
		A int    `vd:"$<=110;msg:'hello'"`
		B string `vd:"len($)>1 && regexp('^\\w*$');msg:'hehe'"`
	}
	t := &TT{
		A: 107,
		B: "abc",
	}

	result, err := Validate(t)
	fmt.Println(result, err)
}
