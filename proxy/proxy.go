package proxy

import (
	"context"
	"errors"
	"fmt"
	"math"
)

var (
	// ErrUnable 代理不可用
	ErrUnable = errors.New("unable to connect proxy")
	// ErrNonEndpoint 没有可用的代理IP
	ErrNonEndpoint = errors.New("no endpoint available")
	// ErrInsufficient 代理的余额不足
	ErrInsufficient = errors.New("insufficient balance")
	// ErrPackageExpired 套餐过期
	ErrPackageExpired = errors.New("package expired")
	// ErrOutException 请求结果异常
	ErrResultException = errors.New("result exception")
)

// Proxy 代理接口
type Proxy interface {
	// 代理初始化
	Init() error

	// 获取代理点
	GetEndpoints(ctx context.Context, n int) ([]*Endpoint, error)

	// 获取代理余额
	GetBalance(ctx context.Context) (*Balance, error)

	// 获取代理信息
	GetJSON() Proxy
}

// 代理余额货币种类，一般根据代理所在的国家或地区决定
type Coin string

const (
	// 人民币
	RMB Coin = "RMB"
)

// INFBalance 无限余额，免费代理的余额
var INFBalance = &Balance{Amount: math.MaxFloat64, Coin: "X"}

// Balance 第三方的余额
type Balance struct {
	// Amount 余额金额
	Amount float64 `json:"amount,omitempty"`

	// Coin 货币类型
	Coin Coin `json:"cain,omitempty"`
}

type Scheme string

const (
	HTTP  Scheme = "http"
	HTTPS Scheme = "https"
	SOCK5 Scheme = "sock5"
)

// Country IP 所在的国家
type Country string

const (
	// China 中国
	China Country = "China"
)

func (c Country) String() string {
	return string(c)
}

// ISP 提供IP的运营商
type ISP string

const (
	DX ISP = "电信"
	LT ISP = "联通"
	YD ISP = "移动"
)

// Endpoint 代理点
type Endpoint struct {
	Scheme Scheme `json:"scheme,omitempty"`

	// Country 代理点所在的国家
	Country Country `json:"country,omitempty"`

	// City 代理点所在的城市
	City string `json:"city,omitempty"`

	// ExpireTime IP 过期时间, 格式为 2006-01-02 15:04:05
	ExpireTime string `json:"expire_time,omitempty"`

	// IP 代理点地址
	IP string `json:"ip,omitempty"`

	// Port 代理点的端口
	Port int `json:"port,omitempty"`

	// ISP 网络运营商
	ISP ISP `json:"isp,omitempty"`

	// 代理点使用次数
	Num int64 `json:"num,omitempty"`
}

func (e *Endpoint) Addr() string {
	return fmt.Sprintf("%s://%s:%d", e.Scheme, e.IP, e.Port)
}

func (e *Endpoint) DeepCopy() *Endpoint {
	out := &Endpoint{}
	*out = *e
	return out
}
