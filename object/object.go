package object

import "errors"

var (
	// ErrObject object 状态错误
	ErrObject = errors.New("object fault")
	// ErrURLExists url 已存在
	ErrURLExists = errors.New("url already exists")
	// ErrOldURL url 已经被访问过
	ErrOldURL = errors.New("url already be visited")
	// ErrNoURL object 没有 url
	ErrNoURL = errors.New("no url in object")
	// ErrGetURL 获取 url 时的错误
	ErrGetURL = errors.New("get object url")
	// ErrSetURL 设置 url 时的错误
	ErrSetURL = errors.New("set object url")
)

type Object interface {
	// 初始化
	Init() error

	// 获取 URL
	GetURL() (*URL, error)

	// 添加 URL
	SetURL(url *URL) error

	// Object 重置，删除内部所有的 url
	Reset()
}

// URL
type URL struct {
	// URL 路径
	Path string `json:"path,omitempty"`

	// URL 所在的 URL 池
	Obj Object `json:"obj,omitempty"`
}
