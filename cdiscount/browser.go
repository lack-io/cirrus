package cdiscount

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"

	"github.com/lack-io/cirrus/internal/log"
	"github.com/lack-io/cirrus/internal/parser"
	"github.com/lack-io/cirrus/proxy"
	"github.com/lack-io/cirrus/storage"
	"github.com/lack-io/cirrus/store"
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
	sub, err := c.storage.Subscribe()
	if err != nil {
		log.Fatalf("订阅 storage 消息错误: %v", err)
	}
START:
	select {
	case <-c.startCh:
	}
	for {
		select {
		case <-c.ctx.Done():
			return
		case u, ok := <-sub.Channel():
			if ok {
				log.Infof("<=== 获取请求路径 %v", u.Path)
				c.goPool.NewTask(func() {
					c.do(u.Path)
				})
			}
		case <-c.pauseCh:
			goto START
		}
	}
}

// StartDaemon implemented daemon.Daemon interfaces
func (c *Cdiscount) StartDaemon(root string) {
	c.storage.Reset()
	_ = c.storage.Push(storage.URL{Path: root})
	c.startCh <- struct{}{}
}

// PauseDaemon implemented daemon.Daemon interfaces
func (c *Cdiscount) PauseDaemon() {
	c.pauseCh <- struct{}{}
}

// do 请求 url
func (c *Cdiscount) do(url string) {
	url, kind := urlParser(url)
	if kind == Unknown {
		log.Infof("目录路径 %s 无效", url)
		c.storage.DelURL(&storage.URL{Path: url})
		return
	}
	c.runTask(url, kind)
}

func (c *Cdiscount) runTask(url string, kind Kind) {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	var err error
	defer func() {
		if err != nil {
			log.Errorf("请求 %s 失败: %v", url, err)
			return
		}
		log.Infof("请求 %s 成功 !!!", url)
		_ = c.storage.DelURL(&storage.URL{Path: url})
	}()

	log.Infof("请求路径 %v", url)
	var doc string
	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector('body')`, &doc, chromedp.ByJSPath),
	}

	var endpoint *proxy.Endpoint
	endpoint, err = c.ProxyPool.GetEndpoint(ctx)
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
	q, err := parser.NewParser(doc)
	if err != nil {
		return
	}

	switch kind {
	case Group:
		for _, node := range q.Each("body", "a") {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					log.Infof("===> 保存请求路径 %v", attr.Val)
					_ = c.storage.SetURL(&storage.URL{Path: attr.Val, Storage: c.storage})
					continue
				}
			}
		}
	case Link:
		for _, node := range q.Each("body", "a") {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					log.Infof("===> 保存请求路径 %v", attr.Val)
					_ = c.storage.SetURL(&storage.URL{Path: attr.Val, Storage: c.storage})
					continue
				}
			}
		}

		// 获取评论信息
		commentTag := strings.TrimSpace(q.Text(".fpTMain .fpDesCol .fpCusto"))
		commentSt := strings.SplitN(commentTag, " ", 2)
		comments, _ := strconv.ParseInt(strings.TrimSpace(commentSt[0]), 10, 64)
		if !(comments > 0) {
			return
		}

		// 获取发货渠道信息
		expressDocs := []string{}
		q.For("#fpShipping .fpShippingMessage", "li .fpShippingText", func(s string, node *html.Node) {
			expressDocs = append(expressDocs, s)
		})
		if len(expressDocs) != 2 {
			return
		}
		if !strings.Contains(expressDocs[1], "Livraison Gratuite") {
			return
		}

		infoDocs := []string{}
		// 获取宝贝信息(特别是宝贝的品牌)
		q.For("#fpContent #descContent table", "tbody tr td", func(s string, node *html.Node) {
			infoDocs = append(infoDocs, strings.TrimSpace(s))
		})
		var index int
		var item string
		for index, item = range infoDocs {
			if item == "Marque" {
				break
			}
		}
		if index < len(infoDocs) && infoDocs[index+1] == "AUCUNE" {
			// 保存符合的宝贝
			good := store.Good{
				URL:       url,
				UID:       urlToID(url),
				Comments:  int(comments),
				Express:   "AUCUNE",
				Timestamp: time.Now().Unix(),
			}
			log.Infof("保存符合要求的宝贝: %v", good.UID)
			err := c.store.AddGood(&good)
			if err != nil {
				log.Errorf("保存宝贝 %s 失败: %v", good.UID, err)
			}
		}
	}
	log.Infof("页面 %s 解析结束!", url)
}

// urlToID 从宝贝的路径提取id
func urlToID(url string) string {
	var id string
	idx := strings.LastIndex(url, "/")
	if idx != -1 {
		id = url[idx+1:]
	}
	if strings.HasSuffix(id, ".html") {
		id = strings.TrimSuffix(id, ".html")
	}
	return id
}

// urlParser 返回处理过的 url 和 url 的类型
func urlParser(url string) (string, Kind) {
	if !strings.HasPrefix(url, prefix) {
		return url, Unknown
	}

	index := strings.LastIndex(url, "/")
	if index == -1 {
		return url, Unknown
	}

	if idx := strings.LastIndex(url, "?"); idx != -1 {
		url = url[:idx]
	}
	if strings.HasPrefix(url, "f") {
		return url, Link
	}

	return url, Group
}
