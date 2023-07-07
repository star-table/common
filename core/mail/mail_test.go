package mail

import (
	"github.com/star-table/common/core/config"
	"testing"
)

func TestAvailability(t *testing.T) {
	t.Logf("start load config")
	config.LoadUnitTestConfig()

	err := SendMail([]string{"ainililia@163.com"}, "hello", "hello")
	if err == nil {
		t.Log("successful")
	} else {
		t.Log(err)
		t.Log("send failed")
	}
}
