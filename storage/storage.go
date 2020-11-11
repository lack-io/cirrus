package storage

import "errors"

var (
	// ErrStorage storage 状态错误
	ErrStorage = errors.New("storage fault")
	// ErrURLExists url 已存在
	ErrURLExists = errors.New("url already exists")
	// ErrOldURL url 已经被访问过
	ErrOldURL = errors.New("url already be visited")
	// ErrNoURL storage 没有 url
	ErrNoURL = errors.New("no url in storage")
	// ErrGetURL 获取 url 时的错误
	ErrGetURL = errors.New("get storage url")
	// ErrSetURL 设置 url 时的错误
	ErrSetURL = errors.New("set storage url")
	// ErrDelURL 删除 url 时的错误
	ErrDelURL = errors.New("del storage url")
)

type Storage interface {
	// 初始化
	Init() error

	// 获取 URL
	GetURL() (*URL, error)

	// 添加 URL
	SetURL(url *URL) error

	// Storage 重置，删除内部所有的 url
	Reset()
}

// URL
type URL struct {
	// URL 路径
	Path string `json:"path,omitempty"`

	// URL 所在的 URL 池
	Storage Storage `json:"storage,omitempty"`
}
