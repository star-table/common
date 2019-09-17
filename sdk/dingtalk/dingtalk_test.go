package dingtalk

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"testing"
)

func TestCreateClient(t *testing.T) {
	t.Logf("start load config")
	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")

	client, _ := GetDingTalkClient("4obevC6UCuxMPFOKKRtCleAzbp9Pz6ft3dHDiiXAEkmD65hs9Sdh5N4fGw2307Hri65huD1IoCapeM8TnE4s8V", "ding8ac2bab2b708b3cc35c2f4657eb6378f")
	t.Log(client)
}
