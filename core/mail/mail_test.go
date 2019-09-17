package mail

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"testing"
)

func TestAvailability(t *testing.T) {
	t.Logf("start load config")
	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-testing\\configs", "application")

	err := SendMail([]string{"ainililia@163.com"}, "hello", "hello")
	if err == nil {
		t.Log("successful")
	} else {
		t.Log("send failed")
	}
}
