package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestBuildSystemErrorInfo(t *testing.T) {
	ssy := BuildSystemErrorInfo(FileNotExist, errors.New("test"))

	fmt.Println(ssy.Message())

	fmt.Println(FileNotExist.Message())

}
