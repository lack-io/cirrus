package net

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	// ErrPostData http request body 数据错误
	ErrPostData = errors.New("bad body data")
	// ErrBadFile 文件数据错误
	ErrBadFile = errors.New("bad file data")
)

type ContentType string

const (
	Json     ContentType = "application/json"
	Form     ContentType = "application/x-www-form-urlencoded"
	FormData ContentType = "multipart/form-data"
)

func (c ContentType) String() string {
	return string(c)
}

// HTTPClient http 客户端，保存 http 请求所需的参数和数据
type HTTPClient struct {
	// http 请求头
	header map[string]string

	// url 参数
	params map[string]string

	// 请求超时时间
	timeout time.Duration

	// 代理地址，格式为 IP:Port
	proxy string

	// http Bearer Token, token 优先
	token string

	// http cookies
	cookies []*http.Cookie

	// Authorization Basic, username 和 password 都不为空时有效
	// token 的优先级更高，token 不为空时选择 Authorization: Bearer
	username string

	password string
}

func Builder() *HTTPClient {
	header := map[string]string{}
	header["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36"
	header["Content-Type"] = Json.String()
	header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	header["Accept-Language"] = "en;q=0.9"

	return &HTTPClient{header: header}
}

func (c *HTTPClient) SetHeader(key, value string) *HTTPClient {
	c.header[key] = value
	return c
}

func (c *HTTPClient) SetParams(params map[string]string) *HTTPClient {
	c.params = params
	return c
}

func (c *HTTPClient) AddCookie(ck *http.Cookie) *HTTPClient {
	c.cookies = append(c.cookies, ck)
	return c
}

func (c *HTTPClient) SetToken(token string) *HTTPClient {
	c.token = token
	return c
}

func (c *HTTPClient) SetBasicAuth(username, password string) *HTTPClient {
	c.username, c.password = username, password
	return c
}

func (c *HTTPClient) SetTimeout(timeout time.Duration) *HTTPClient {
	c.timeout = timeout
	return c
}

func (c *HTTPClient) SetContentType(ct ContentType) *HTTPClient {
	c.header["Content-Type"] = ct.String()
	return c
}

// Do 自定义请求
//	method: GET POST PATCH PUT DELETE HEAD
//	path: 请求 url
//	params: url 参数
//	data: http request body
//	files: 待上传的文件
func (c *HTTPClient) Do(method, path string, data, files map[string]string) ([]byte, error) {

	hc := http.Client{}

	// params 不为空时，拼接 path
	if c.params != nil {
		query := []string{}
		for k, v := range c.params {
			query = append(query, k+"="+url.PathEscape(v))
		}
		path = path + "?" + strings.Join(query, "&")
	}

	req := &http.Request{
		Method: method,
		URL:    &url.URL{Path: path},
	}

	// 设置请求头
	for k, v := range c.header {
		req.Header.Set(k, v)
	}

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	if c.cookies != nil {
		for _, ck := range c.cookies {
			req.AddCookie(ck)
		}
	}

	// data 不为 nil 时，设置 http request Body
	if data != nil {
		bodyBuf := &bytes.Buffer{}
		bw := multipart.NewWriter(bodyBuf)

		for k, v := range data {
			_ = bw.WriteField(k, v)
		}

		if files != nil {
			for f, p := range files {
				fw, err := bw.CreateFormFile(f, f)
				if err != nil {
					return nil, fmt.Errorf("%w: %v", ErrPostData, err)
				}

				//打开文件句柄操作
				file, err := os.Open(p)
				if err != nil {
					return nil, fmt.Errorf("%w: %v", ErrBadFile, err)
				}

				_, err = io.Copy(fw, file)
				if err != nil {
					return nil, fmt.Errorf("%w: %v", ErrBadFile, err)
				}
				_ = file.Close()
			}
		}

		err := bw.Close()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPostData, err)
		}

		c.header["Content-Type"] = bw.FormDataContentType()

		// io.Read => io.ReadCloser
		req.Body = ioutil.NopCloser(bodyBuf)
	}

	// 开始请求
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Get 发起一次 GET 请求
func (c *HTTPClient) Get(url string) ([]byte, error) {
	return c.Do("GET", url, nil, nil)
}

// PostFile 上传文件
//	path: 请求 url
//	data: http request body
//	files: 待上传的文件
func (c *HTTPClient) PostFile(url string, data, files map[string]string) ([]byte, error) {
	return c.Do("POST", url, data, files)
}

// Post 发起一次 POST 请求
//	path: 请求 url
//	data: http request body
func (c *HTTPClient) Post(url string, data map[string]string) ([]byte, error) {
	return c.Do("POST", url, data, nil)
}

// Patch 发起一次 PATCH 请求
//	path: 请求 url
//	data: http request body
func (c *HTTPClient) Patch(url string, data map[string]string) ([]byte, error) {
	return c.Do("PATCH", url, data, nil)
}

// Put 发起一次 PUT 请求
//	path: 请求 url
//	data: http request body
func (c *HTTPClient) Put(url string, data map[string]string) ([]byte, error) {
	return c.Do("PUT", url, data, nil)
}

// Delete 发起一次 DELETE 请求
//	path: 请求 url
//	data: http request body
func (c *HTTPClient) DELETE(url string, data map[string]string) ([]byte, error) {
	return c.Do("DELETE", url, data, nil)
}
