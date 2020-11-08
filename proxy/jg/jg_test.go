package jg

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/xingyys/cirrus/config"
)

var jg = &JG{}

func TestGetJGProxy(t *testing.T) {
	jg, _ = GetJGProxy(context.Background(), &config.ProxyJG{
		Neek:          "29394",
		APIAppKey:     "ec28da82d45f195e3d962c86dce693f4",
		BalanceAppKey: "ec36c20228d8a97ff5dc56273dda535d",
	})
}

func TestJG_Init(t *testing.T) {
	err := jg.Init()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, jg.Whites, nil)
	assert.NotEqual(t, jg.LocalIP, "")

	time.Sleep(time.Second * 5)

	t.Log(jg.Whites)
	t.Log(jg.LocalIP)
}

func TestJG_GetBalance(t *testing.T) {
	b, err := jg.GetBalance(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)
}

func TestJG_GetPackageBalance(t *testing.T) {
	b, err := jg.GetPackageBalance(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)
}

func TestJG_GetEndpoint(t *testing.T) {
	eps, err := jg.GetEndpoints(context.TODO(), 3)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(eps)
	t.Log(string(data))
}
