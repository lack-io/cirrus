package parser

import (
	"context"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"

	"github.com/lack-io/cirrus/internal/client"
	"github.com/lack-io/cirrus/internal/net"
)

func newClient() *client.Client {
	opts := client.Option{
		Headless:                false,
		BlinkSettings:           "imagesEnabled=false",
		UserAgent:               net.UserAgent,
		IgnoreCertificateErrors: true,
		WindowsHigh:             400,
		WindowsWith:             400,
	}

	return client.NewClient(context.TODO(), opts)
}

func TestQuery_Each(t *testing.T) {
	cli := newClient()

	var html string
	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector('body')`, &html, chromedp.ByJSPath),
	}

	err := cli.NewTask().
		ExecOption(chromedp.ProxyServer("http://127.0.0.1:7890")).
		Actions(actions...).
		Do(context.TODO(), `https://www.cdiscount.com/`)
	if err != nil {
		t.Fatal(err)
	}

	q, err := NewParser(html)
	if err != nil {
		t.Fatal(err)
	}

	for _, node := range q.Each("body", "a") {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				t.Log(attr.Val)
			}
		}
	}
}

func TestParser_Text(t *testing.T) {
	cli := newClient()

	var html string
	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector('body')`, &html, chromedp.ByJSPath),
	}

	err := cli.NewTask().
		ExecOption(chromedp.ProxyServer("http://127.0.0.1:7890")).
		Actions(actions...).
		Do(context.TODO(), `https://www.cdiscount.com/informatique/gaming/manette-sans-fil-pour-nintendo-switch-bluetooth-m/f-107140308-auc6954248714547.html?idOffre=692865572`)
	if err != nil {
		t.Fatal(err)
	}

	q, err := NewParser(html)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(q.Text("#fpContent #descContent table tbody"))
}

func TestParser_Html(t *testing.T) {
	cli := newClient()

	var html string
	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector('body')`, &html, chromedp.ByJSPath),
	}

	err := cli.NewTask().
		ExecOption(chromedp.ProxyServer("http://127.0.0.1:7890")).
		Actions(actions...).
		Do(context.TODO(), `https://www.cdiscount.com/informatique/gaming/manette-sans-fil-pour-nintendo-switch-bluetooth-m/f-107140308-auc6954248714547.html?idOffre=692865572`)
	if err != nil {
		t.Fatal(err)
	}

	q, err := NewParser(html)
	if err != nil {
		t.Fatal(err)
	}

	hl, err := q.Html("#fpContent #descContent table tbody")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hl)
}

func TestParser_For(t *testing.T) {
	cli := newClient()

	var out string
	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector('body')`, &out, chromedp.ByJSPath),
	}

	err := cli.NewTask().
		ExecOption(chromedp.ProxyServer("http://127.0.0.1:7890")).
		Actions(actions...).
		Do(context.TODO(), `https://www.cdiscount.com/informatique/gaming/manette-sans-fil-pour-nintendo-switch-bluetooth-m/f-107140308-auc6954248714547.html`)
	if err != nil {
		t.Fatal(err)
	}

	q, err := NewParser(out)
	if err != nil {
		t.Fatal(err)
	}

	// 获取评论信息
	t.Log(strings.TrimSpace(q.Text(".fpTMain .fpDesCol .fpCusto")))

	// 获取发货渠道信息
	q.For("#fpShipping .fpShippingMessage", "li .fpShippingText", func(s string, node *html.Node) {
		t.Logf("text = %v, node = %v\n", s, node)
	})

	// 获取宝贝信息(特别是宝贝的品牌)
	//q.For("#fpContent #descContent table", "tbody tr td", func(s string, node *html.Node) {
	//	t.Logf("text = %v, node = %v\n", s, node)
	//})
}