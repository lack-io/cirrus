package cdiscount

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"

	"github.com/xingyys/cirrus/config"
	"github.com/xingyys/cirrus/internal/log"
	"github.com/xingyys/cirrus/proxy"
	"github.com/xingyys/cirrus/proxy/jg"
)

// Pool 代理IP池，负责管理代理Pool
type Pool struct {
	ctx    context.Context
	cancel context.CancelFunc

	// opts 代理参数
	opts *config.Proxy

	// 可用的代理集合
	pp proxy.Proxy

	// 代理IP池
	endpoints *sll.List
	// endpoints 读写锁
	elock *sync.RWMutex

	// 代理失效信号，每次检测到代理失效时，传入值
	expiredCh chan struct{}
	// 获取代理时的错误信息
	proxyErrCh chan error
}

// NewPool 新建 Pool
func NewPool(ctx context.Context, opts *config.Proxy) (*Pool, error) {
	ctx, cancel := context.WithCancel(ctx)
	p := &Pool{
		ctx:        ctx,
		cancel:     cancel,
		opts:       opts,
		endpoints:  sll.New(),
		elock:      &sync.RWMutex{},
		expiredCh:  make(chan struct{}, 1),
		proxyErrCh: make(chan error, 1),
	}

	if opts.Enable {
		switch opts.Agent {
		case config.JG:
			if opts.JG == nil {
				return nil, fmt.Errorf("config is nil")
			}
			jgProxy, err := jg.GetJGProxy(ctx, opts.JG)
			if err != nil {
				return nil, err
			}
			if err := jgProxy.Init(); err != nil {
				return nil, err
			}
			p.pp = jgProxy
			log.Info("init [jiguang] proxy successfully")
		}
	}

	endpoints, err := p.pp.GetEndpoints(ctx, p.opts.Size)
	if err != nil {
		return nil, err
	}
	for _, e := range endpoints {
		p.endpoints.Add(e)
	}

	go p.process()

	return p, nil
}

// GetEndpoint 从代理池中随机获取一个代理点，如果获取的代理点即将过期，则删除并从新获取
func (p *Pool) GetEndpoint(ctx context.Context) (*proxy.Endpoint, error) {
	p.elock.RLock()
	size := p.endpoints.Size()
	p.elock.RUnlock()
	if size == 0 {
		return nil, fmt.Errorf("%w: pool is empty", proxy.ErrNonEndpoint)
	}

	var endpoint *proxy.Endpoint
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-p.proxyErrCh:
			return nil, err
		default:

		}
		p.elock.RLock()
		rand.Seed(time.Now().UnixNano())
		n := p.endpoints.Size()
		rd := rand.Intn(n)
		v, ok := p.endpoints.Get(rd)
		p.elock.RUnlock()
		if ok {
			endpoint = v.(*proxy.Endpoint)
			t, _ := time.Parse("2006-01-02 15:03:04", endpoint.ExpireTime)
			if t.Sub(time.Now()).Seconds() > 30 {
				break
			}
			log.Infof("proxy endpoint [%s:%d] expired", endpoint.IP, endpoint.Port)
			p.elock.Lock()
			p.endpoints.Remove(rd)
			p.elock.Unlock()
			p.expiredCh <- struct{}{}
		}
	}

	return endpoint, nil
}

// process 代理池后台任务，等待某个代理节点过期的信号，接收到信号时，从第三方代理获取一个新的代理节点
func (p *Pool) process() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-p.expiredCh:
			log.Infof("get new proxy endpoint")
			endpoint, err := p.pp.GetEndpoints(p.ctx, 1)
			if err != nil {
				p.proxyErrCh <- err
			} else {
				p.elock.Lock()
				p.endpoints.Add(endpoint)
				p.elock.Unlock()
			}
		}
	}
}

func (p *Pool) Endpoints() ([]*proxy.Endpoint, error) {
	p.elock.RLock()
	defer p.elock.RUnlock()

	endpoints := make([]*proxy.Endpoint, 0)
	values := p.endpoints.Values()
	for _, item := range values {
		endpoints = append(endpoints, item.(*proxy.Endpoint))
	}
	return endpoints, nil
}

func (p *Pool) Close() {
	p.cancel()
}
