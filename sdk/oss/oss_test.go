package oss

import (
	"bufio"
	"fmt"
	"github.com/galaxy-book/common/core/config"
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

func TestGetObjectUrl(t *testing.T) {

	config.LoadUnitTestConfig()
	signUrl, err := GetObjectUrl("org_1325/project_1525/issue_7127/resource/2019/12/12/91d6d6461c6e4176a77aa486dde75f651576153986043.xlsx",  6000)
	t.Log(err)
	t.Log(signUrl)
}