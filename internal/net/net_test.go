package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicIP(t *testing.T) {
	pip, err := GetPublicIP()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, pip.IP, "")

	t.Log(pip)
}

func TestParseProxy(t *testing.T) {
	proxy := "127.0.0.1:3000"

	ip, port, err := ParseProxy(proxy)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ip, "127.0.0.1")
	assert.Equal(t, port, "3000")
}

func TestCheckProxyIP(t *testing.T) {
	proxy := "223.243.72.8:45351"

	ok, err := CheckProxyIP(proxy)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, ok)
}