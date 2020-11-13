package redis

import (
	"context"
	"testing"

	"github.com/lack-io/cirrus/config"
	"github.com/lack-io/cirrus/storage"
)

var ob = &Redis{}

const (
	addr = "192.168.3.111:6379"
	username = ""
	password = ""
)


func TestNewRedis(t *testing.T) {
	var err error
	ctx := context.Background()
	cfg := &config.StorageRedis{}
	ob = NewRedis(ctx, cfg)
	err = ob.Init()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedis_Push(t *testing.T) {
	url := storage.URL{Path: "https://www.google.com"}

	err := ob.Push(url)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedis_Reset(t *testing.T) {
	ob.Reset()
}

