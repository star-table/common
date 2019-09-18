package oss

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"testing"
)

func TestPostPolicy(t *testing.T) {

	config.LoadUnitTestConfig()
	pp := PostPolicy("project", 1000*60*5, 0)
	t.Log(json.ToJson(pp))
}
