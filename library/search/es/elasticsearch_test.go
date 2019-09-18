package es

import (
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"testing"
)

func TestGetESClient(t *testing.T) {
	config.LoadUnitTestConfig()

	client, err := GetESClient()
	fmt.Println(err)
	fmt.Printf("%v \n", client)
}
