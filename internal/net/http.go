package net

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	urlpkg "net/url"
	"strings"
	"time"
)

const (
	JSON = "application/json"
	FORM = "application/x-www-form-urlencoded"
	MUL  = "multipart/form-data"
)

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
	header["Content-Type"] = JSON
	header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	header["Accept-Language"] = "en;q=0.9"

	return &HTTPClient{header: header, params: map[string]string{}}
}

func (c *HTTPClient) SetHeader(key, value string) *HTTPClient {
	c.header[key] = value
	return c
}

func (c *HTTPClient) AddParams(params map[string]string) *HTTPClient {
	for k, v := range params {
		c.params[k] = v
	}
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

func (c *HTTPClient) SetContentType(ct string) *HTTPClient {
	c.header["Content-Type"] = ct
	return c
}

// Do 自定义请求
//	method: GET POST PATCH PUT DELETE HEAD
//	path: 请求 url
//	params: url 参数
//	data: http request body
//	files: 待上传的文件
func (c *HTTPClient) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {

	hc := http.Client{}

	// params 不为空时，拼接 path
	if c.params != nil {
		query := []string{}
		for k, v := range c.params {
			query = append(query, k+"="+urlpkg.PathEscape(v))
		}
		url = url + "?" + strings.Join(query, "&")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
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

	// 开始请求
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Get 发起一次 GET 请求
func (c *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	return c.Do(ctx, "GET", url, nil)
}

// Post 发起一次 POST 请求
func (c *HTTPClient) Post(ctx context.Context, url string, body io.Reader) ([]byte, error) {
	return c.Do(ctx, "POST", url, body)
}

// Patch 发起一次 PATCH 请求
func (c *HTTPClient) Patch(ctx context.Context, url string, body io.Reader) ([]byte, error) {
	return c.Do(ctx, "PATCH", url, body)
}

// Put 发起一次 PUT 请求
func (c *HTTPClient) Put(ctx context.Context, url string, body io.Reader) ([]byte, error) {
	return c.Do(ctx, "PUT", url, body)
}

// Delete 发起一次 DELETE 请求
func (c *HTTPClient) DELETE(ctx context.Context, url string, body io.Reader) ([]byte, error) {
	return c.Do(ctx, "DELETE", url, body)
}
