// 极光爬虫代理
package jg

import (
	"context"
	"time"

	"github.com/xingyys/cirrus/proxy"
)

// While JG 报名单信息
type White struct {
	// 白名单IP
	IP string `json:"ip,omitempty"`

	// 更新时间，格式为时间戳，单位为秒
	UpdateAt string `json:"updateAt,omitempty"`
}

type JG struct {
	// ctx 控制 JG 的停止
	ctx context.Context

	Name string `json:"name,omitempty"`

	// 代理点容量
	Size int `json:"size,omitempty"`

	// 极光代理余额信息
	Balance *proxy.Balance `json:"balance,omitempty"`

	// 白名单信息
	Whites []White `json:"whites,omitempty"`
}

func GetJGProxy(ctx context.Context, size int, timeout time.Duration) *JG {

	return &JG{ctx: ctx, Name: "jiguang", Size: size}
}

// Init implement proxy.Proxy
func (j *JG) Init() error {
	return nil
}

// GetEndpoint implement proxy.Proxy
func (j *JG) GetEndpoint() (*proxy.Endpoint, error) {
	// TODO: GetEndpoint
	return nil, nil
}

// GetBalance implement proxy.Proxy
func (j *JG) GetBalance() *proxy.Balance {
	return j.Balance
}
