package client

import (
	"context"
	"sync"
	"testing"

	"github.com/chromedp/chromedp"

	"github.com/xingyys/cirrus/internal/net"
)

var cli *Client = &Client{}

func newClient() {
	opts := Option{
		Headless:                false,
		BlinkSettings:           "imagesEnabled=false",
		UserAgent:               net.UserAgent,
		IgnoreCertificateErrors: true,
		WindowsHigh:             400,
		WindowsWith:             400,
	}

	cli = NewClient(context.TODO(), opts)
}

func TestClient_NewTask(t *testing.T) {
	newClient()

	var html string

	actions := []chromedp.Action{
		chromedp.WaitVisible(`#page`, chromedp.ByID),
		chromedp.OuterHTML(`document.querySelector("#paContent")`, &html, chromedp.ByJSPath),
	}

	err := cli.NewTask().ExecOption().Actions(actions...).Do(context.TODO(), `https://www.cdiscount.com/`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(html)
}

func TestClient_NewProxyTask(t *testing.T) {
	newClient()

	var html string

	option := chromedp.ProxyServer("http://127.0.0.1:7890")

	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector("body")`, &html, chromedp.ByJSPath),
	}

	err := cli.NewTask().ExecOption(option).Actions(actions...).Do(context.TODO(), `https://www.cdiscount.com/`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(html)
}

func TestClient_NewInvalidProxy(t *testing.T) {
	newClient()

	var html string

	option := chromedp.ProxyServer("http://127.0.0.1:78911")

	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector("body")`, &html, chromedp.ByJSPath),
	}

	err := cli.NewTask().ExecOption(option).Actions(actions...).Do(context.TODO(), `https://www1.baidu.com/`)
	t.Log(err)
	t.Log(html)
}

func TestClient_NewTask_QueryLink(t *testing.T) {
	newClient()

	var html string
	actions := []chromedp.Action{
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`document.querySelector('body')`, &html, chromedp.ByJSPath),
	}

	err := cli.NewTask().
		ExecOption(chromedp.ProxyServer("http://127.0.0.1:7890")).
		Actions(actions...).
		Do(context.TODO(), `https://www.cdiscount.com/high-tech/sono-dj/karaoke-600w-2-micros-fil-et-sans-fil-enceinte-s/f-1063919-boo7112132699946.html`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(html)
}

func TestTask_DoubleDo(t *testing.T) {
	newClient()

	links := []string{`https://www.baidu.com`, `https://www.qq.com`}
	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var html string
			actions := []chromedp.Action{
				chromedp.WaitReady(`body`, chromedp.ByQuery),
				chromedp.OuterHTML(`document.querySelector('body')`, &html, chromedp.ByJSPath),
			}
			err := cli.NewTask().
				ExecOption(chromedp.ProxyServer("http://127.0.0.1:7890")).
				Actions(actions...).
				Do(context.TODO(), links[i])
			if err != nil {
				t.Fatal(err)
			}
			t.Log(html)
		}(i)
	}

	wg.Wait()
}
