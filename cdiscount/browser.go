package cdiscount

import (
	"context"
	urlpkg "net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/lack-io/cirrus/internal/log"
	"github.com/lack-io/cirrus/internal/parser"
	"github.com/lack-io/cirrus/proxy"
)

const (
	prefix = "https://www.cdiscount.com"
)

type Kind string

const (
	Unknown Kind = "unknown"
	Group   Kind = "group"
	Link    Kind = "link"
)

func (c *Cdiscount) daemon() {
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-c.ctx.Done():
			close(c.connects)
			timer.Stop()
			return
		case <-timer.C:
			u, _ := c.storage.GetURL()
			if u != nil {
				log.Infof("<=== 获取请求路径 %v", u.Storage, u.Path)
				c.connects <- u.Path
			}
		case url, ok := <-c.connects:
			if !ok {
				break
			}

			go c.do(url)
		}
	}
}

// do 请求 url
func (c *Cdiscount) do(url string) {
	url, kind := urlParser(url)
	switch kind {
	case Unknown:
		log.Infof("目录路径 %s 无效", url)
		return
	case Group:
		c.doGroup(url)
	case Link:

	}

}

func (c *Cdiscount) doGroup(url string) {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	var err error
	defer func() {
		if err != nil {
			log.Errorf("请求 %s 失败: %v", url, err)
			c.connects <- url
			return
		}
		log.Infof("请求 %s 成功 !!!", url)
	}()

	log.Infof("请求路径 %v", url)
	var html string
	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector('body')`, &html, chromedp.ByJSPath),
	}

	var endpoint *proxy.Endpoint
	endpoint, err = c.Pool.GetEndpoint(ctx)
	if err != nil {
		return
	}

	log.Infof("获取代理节点 %v", endpoint.Addr())
	err = c.cli.NewTask().
		ExecOption(chromedp.ProxyServer(endpoint.Addr())).
		Actions(actions...).
		Do(ctx, url)
	if err != nil {
		return
	}

	log.Infof("开始解析 %v 页面...", url)
	q, err := parser.NewParser(html)
	if err != nil {
		return
	}

	for _, node := range q.Each("body", "a") {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				_, kind := urlParser(attr.Val)
				if kind != Unknown {
					log.Infof("===> 保存请求路径 %v", attr.Val)
				}
			}
		}
	}
}

// urlParser 返回处理过的 url 和 url 的类型
func urlParser(url string) (string, Kind) {
	URL, err := urlpkg.Parse(url)
	if err != nil {
		return "", Unknown
	}

	path := URL.Path
	if !strings.HasPrefix(path, prefix) {
		return url, Unknown
	}

	index := strings.LastIndex(path, "/")
	if index == -1 {
		return url, Unknown
	}

	if strings.HasPrefix(path, "f") {
		return path, Link
	}

	return path, Group
}
