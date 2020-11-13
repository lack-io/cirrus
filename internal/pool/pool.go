package pool

import "context"

type Pool struct {
	ctx  context.Context
	work chan func()
	sem  chan struct{}
}

func New(ctx context.Context, size int) *Pool {
	return &Pool{
		ctx:  ctx,
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}

func (p *Pool) NewTask(fn func()) {
	select {
	case p.work <- fn:
	case p.sem <- struct{}{}:
		go p.worker(fn)
	}
}

func (p *Pool) worker(fn func()) {
	defer func() {
		<-p.sem
	}()
	for {
		fn()
		select {
		case <-p.ctx.Done():
			return
		case fn = <-p.work:
		}

	}
}
