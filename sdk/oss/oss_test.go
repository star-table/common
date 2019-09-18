package oss

import (
	"bufio"
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"os"
	"testing"
)

func TestSignUrlWithStream(t *testing.T) {

	keyname := "test/hello5.jpg"
	config.LoadUnitTestConfig()

	signUrl, _ := SignUrlWithStream(keyname, "png", 60)

	fmt.Println(signUrl)

	//f, _ := os.Open("C:\\Users\\admin\\Desktop\\51980670.jpg")
	f, _ := os.Open("C:\\Users\\admin\\Desktop\\dingding-test.png")

	buf := bufio.NewReader(f)

	err := PutObjectWithURL(signUrl, "png", buf)
	if err != nil {
		t.Log(err)
	}

}
