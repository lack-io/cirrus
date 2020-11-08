package config

import (
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	once = &sync.Once{}
	conf = &Config{}
)

type Kind string

const (
	Redis Kind = "redis"
)

// Get 获取全局 config
func Get() *Config {
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

type Config struct {
	Web *Web `toml:"web"`

	Storage *Storage `toml:"object"`

	Store *Store `toml:"store"`

	Client *Client `toml:"Client"`

	Proxy *Proxy `toml:"proxy"`

	Logger *Logger `toml:"logger"`
}

// Web 模块配置
type Web struct {
	// Web 服务绑定地址
	Binding string `toml:"binding"`

	// Web 服务端口
	Port int `toml:"port"`
}

// Storage 模块配置
type Storage struct {
	// URL 存储方式
	Kind Kind `toml:"storage"`

	// Redis 配置，Storage=Redis 时有效
	Redis *StorageRedis `toml:"redis"`
}

// Storage 模块 redis 配置
type StorageRedis struct {
	// Redis 地址
	Addr string `toml:"addr"`

	// Redis 用户名
	Username string `toml:"username"`

	// Redis 密码
	Password string `toml:"password"`

	// Redis 连接数
	Pools int `toml:"pools"`
}

type StoreDB string

const (
	Sqlite StoreDB = "sqlite"
)

type Store struct {
	DB StoreDB `toml:"db"`

	Sqlite *DBSqlite `toml:"sqlite"`
}

type DBSqlite struct {
	Name string `toml:"name"`
}

// Client 模块配置
type Client struct {
	// Headless 是否隐藏 chrome
	Headless bool `toml:"headless"`

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
	Enable bool `toml:"enable"`
	// IP代理商
	Agent Agent `toml:"agents"`

	Size int `toml:"size"`

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

type Logger struct {
	Filename string `toml:"filename"`

	MaxSize int `toml:"maxsize"`

	MaxAge int `toml:"maxage"`

	MaxBackups int `toml:"maxbackups"`

	LocalTime bool `toml:"localtime"`

	Compress bool `toml:"compress"`
}
