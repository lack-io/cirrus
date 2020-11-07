package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := Init("../cirrus-test.toml")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	conf := Get()

	assert.NotEqual(t, conf.Web, nil)
	assert.Equal(t, conf.Web.Binding, "127.0.0.1")
	assert.Equal(t, conf.Web.Port, 4455)
	assert.NotEqual(t, conf.Storage, nil)
	assert.Equal(t, conf.Storage.Kind, Redis)
	assert.NotEqual(t, conf.Storage.Redis, nil)
	assert.Equal(t, conf.Storage.Redis.Addr, "192.168.3.111:6379")
	assert.NotEqual(t, conf.Client, nil)
	assert.NotEqual(t, conf.Proxy, nil)
}
