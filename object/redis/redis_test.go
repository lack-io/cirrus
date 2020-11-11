package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/lack-io/cirrus/object"
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
	ob = NewRedis(ctx, addr, username, password, 10)
	err = ob.Init()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedis_SetURL(t *testing.T) {
	url := &object.URL{Path: "https://www.google.com"}

	err := ob.SetURL(url)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedis_GetURL(t *testing.T) {
	u, err := ob.GetURL()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(u)
}

func TestRedis_Reset(t *testing.T) {
	ob.Reset()
}

func BenchmarkRedis_SetURL(b *testing.B) {
	ctx := context.Background()
	ob = NewRedis(ctx, addr, username, password, 10)
	_ = ob.Init()
	for i := 0; i < b.N; i++ {
		u := object.URL{Path: fmt.Sprintf("https://a%d", i)}
		_ = ob.SetURL(&u)
	}
}

func BenchmarkRedis_GetURL(b *testing.B) {
	ctx := context.Background()
	ob = NewRedis(ctx, addr, username, password, 10)
	_ = ob.Init()
	for i := 0; i < b.N; i++ {
		u, err := ob.GetURL()
		b.Log(u, err)
	}
}