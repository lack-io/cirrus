package config

import (
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	once = &sync.Once{}
	conf = &config{}
)

type Storage string

const (
	Redis Storage = "redis"
)

// Get 获取全局 config
func Get() *config {
	return conf
}

// Init 初始化配置文件
func Init(path string) error {
	var err error
	once.Do(func() {
		_, err = toml.DecodeFile(path, conf)
	})

	return err
}

type config struct {
	Web *Web `toml:"web"`

	Object *Object `toml:"object"`

	Client *Client `toml:"Client"`

	Proxy *Proxy `toml:"proxy"`
}

// Web 模块配置
type Web struct {
	// Web 服务绑定地址
	Binding string `toml:"binding"`

	// Web 服务端口
	Port int `toml:"port"`
}

// Object 模块配置
type Object struct {
	// URL 存储方式
	Storage Storage `toml:"storage"`

	// Redis 配置，Storage=Redis 时有效
	Redis *ObjectRedis `toml:"redis"`
}

// Object 模块 redis 配置
type ObjectRedis struct {
	// Redis 地址
	Addr string `toml:"addr"`

	// Redis 用户名
	Username string `toml:"username"`

	// Redis 密码
	Password string `toml:"password"`

	// Redis 连接数
	Pools int `toml:"pools"`
}

// Client 模块配置
type Client struct {
	// 忽略图片加载
	SkipImage bool `toml:"skip_image"`

	// 并发连接数
	Connections int `toml:"connections"`
}

type Agent string

const (
	JG Agent = "jg"
)

// Proxy 模块配置
type Proxy struct {
	// IP代理商
	Agents []Agent `toml:"agents"`

	// 极光代理配置，Agents 包含 JG 时有效
	JG *ProxyJG `toml:"jg"`
}

// jg 代理配置
type ProxyJG struct {
	// JG 账号的 neek 信息
	Neek string `toml:"neek"`

	// JG 账号的 代理接口 appkey 信息
	APIAppKey string `toml:"api_appkey"`

	// JG 余额接口 appkey 信息
	BalanceAppKey string `toml:"balance_appkey"`
}
