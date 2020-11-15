package parser

import (
	"context"
	"strconv"
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
	ctx := context.TODO()

	var html string
	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector('body')`, &html, chromedp.ByJSPath),
	}

	err := cli.NewTask().
		//ExecOption(chromedp.ProxyServer("http://127.0.0.1:7890")).
		Actions(actions...).
		Do(ctx, `https://www.cdiscount.com/`)
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
		//ExecOption(chromedp.ProxyServer("http://127.0.0.1:10810/pac/?t=20201108170517704")).
		Actions(actions...).
		Do(context.TODO(), `https://www.cdiscount.com/jardin/entretenir-les-plantes/style-de-charge-usb-de-montre-bracelet-usb-anti-mo/f-1630203-auc2008487052282.html`)
	if err != nil {
		t.Fatal(err)
	}

	q, err := NewParser(out)
	if err != nil {
		t.Fatal(err)
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
		t.Log("OK")
	}
}

func TestUrl(t *testing.T) {
	url := `https://www.cdiscount.com/jardin/entretenir-les-plantes/style-de-charge-usb-de-montre-bracelet-usb-anti-mo/f-1630203-auc2008487052282.html?idOffre=600160694#cm_rr=FP:10114852:SW:CAR&sw=33fb8bb4f9ca038ed912313145ebcc9782e70906a099c6e99823b73d4bef41f554e28794ac15ef8eb2dcd0cde149051732536e28c6059eb25ecd410de6da8d8be3fbf3db06069559057929a4ce7215fdc82a5accb2d55fcb3563d92962cfc83876f79126271c84c69db2da98d01ed497ba34c64757a620a624fec2ba884d7105`


	index := strings.LastIndex(url, ".html")
	if index == -1 {
		t.Fail()
	}

	url = url[:index+5]
	t.Log(url)
}