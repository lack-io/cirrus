package redis

import (
	"context"
	"testing"

	"github.com/lack-io/cirrus/config"
	"github.com/lack-io/cirrus/storage"
)

var ob = &Redis{}

const (
	addr     = "192.168.10.10:6379"
	username = ""
	password = ""
)

func newRedis() error {
	var err error
	ctx := context.Background()
	cfg := &config.StorageRedis{
		Addr: addr,
		Pools: 3,
	}
	ob = NewRedis(ctx, cfg)
	err = ob.Init()
	return err
}

func TestRedis_Push(t *testing.T) {
	newRedis()
	url := storage.URL{Path: "https://www.google.com"}

	err := ob.Push(url)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedis_Persist(t *testing.T) {
	newRedis()
	url := storage.URL{Path: "https://www.google.com"}

	err := ob.Persist(url)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedis_Reset(t *testing.T) {
	ob.Reset()
}
