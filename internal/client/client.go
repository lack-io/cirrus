package client

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

type Client struct {
	ctx context.Context

	option Option

	opts []chromedp.ExecAllocatorOption
}

func NewClient(ctx context.Context, option Option) *Client {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", option.Headless),
		chromedp.Flag("blink-settings", option.BlinkSettings),
		chromedp.UserAgent(option.UserAgent),
		chromedp.Flag("ignore-certificate-errors", option.IgnoreCertificateErrors),
		chromedp.WindowSize(option.WindowsHigh, option.WindowsWith),
	)

	c := &Client{
		ctx:    ctx,
		option: option,
		opts:   opts,
	}

	return c
}

// NewTask 启动一起请求任务
func (c *Client) NewTask() *Task {
	return newTask(c)
}

type Task struct {
	cli *Client

	execOptions []chromedp.ExecAllocatorOption

	actions []chromedp.Action
}

func newTask(cli *Client) *Task {
	return &Task{
		cli:         cli,
		execOptions: []chromedp.ExecAllocatorOption{},
		actions:     []chromedp.Action{},
	}
}

func (t *Task) ExecOption(opts ...chromedp.ExecAllocatorOption) *Task {
	t.execOptions = append(t.execOptions, opts...)
	return t
}

func (t *Task) Actions(actions ...chromedp.Action) *Task {
	t.actions = append(t.actions, actions...)
	return t
}

func (t *Task) Do(ctx context.Context, urlstr string) error {
	opts := append(t.cli.opts, t.execOptions...)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	actions := []chromedp.Action{chromedp.Navigate(urlstr)}
	actions = append(actions, t.actions...)

	return chromedp.Run(taskCtx, actions...)
}
