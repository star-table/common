package emt

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnect(t *testing.T) {

	config.LoadUnitTestConfig()

	client, err := GetClient()
	assert.Equal(t, err, nil)
	assert.Equal(t, client.IsConnected(), true)
}

func TestGenerateKey(t *testing.T) {
	config.LoadUnitTestConfig()

	key, err := GenerateKey("nico/#/", "rw", 10000)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, key, "")
	t.Log(key)
}