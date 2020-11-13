package cdiscount

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lack-io/cirrus/config"
	"github.com/lack-io/cirrus/controller"
	"github.com/lack-io/cirrus/internal/client"
	"github.com/lack-io/cirrus/internal/log"
	"github.com/lack-io/cirrus/internal/net"
	"github.com/lack-io/cirrus/internal/pool"
	"github.com/lack-io/cirrus/storage"
	"github.com/lack-io/cirrus/storage/redis"
	"github.com/lack-io/cirrus/store"
)

type Cdiscount struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg *config.Config

	store *store.Store

	ProxyPool *Pool

	cli *client.Client

	Serve *http.Server

	storage storage.Storage

	goPool *pool.Pool

	startCh chan struct{}
	pauseCh chan struct{}
}

func NewCdiscount(cfg *config.Config) (*Cdiscount, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cds := &Cdiscount{
		ctx:     ctx,
		cancel:  cancel,
		cfg:     cfg,
		goPool:  pool.New(ctx, cfg.Client.Connections),
		startCh: make(chan struct{}, 1),
		pauseCh: make(chan struct{}, 1),
	}

	if err := cds.initLogger(); err != nil {
		return nil, err
	}

	log.Info("init data store")
	if err := cds.initStore(); err != nil {
		return nil, err
	}
	log.Info("init data store [ok]")

	log.Info("init storage")
	if err := cds.initStorage(); err != nil {
		return nil, err
	}
	log.Info("init storage [ok]")

	log.Info("init proxy pool")
	if err := cds.initPool(); err != nil {
		return nil, err
	}
	log.Info("init proxy pool")

	log.Info("init chrome client")
	if err := cds.initClient(); err != nil {
		return nil, err
	}
	log.Info("init proxy pool [ok]")

	log.Info("init web server [ok]")
	cds.initServe()

	return cds, nil
}

func (c *Cdiscount) initLogger() error {
	err := log.Init(c.cfg.Logger)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cdiscount) initStore() error {
	s, err := store.NewStore(c.cfg.Store)
	if err != nil {
		return err
	}
	c.store = s
	return nil
}

func (c *Cdiscount) initStorage() error {
	var err error
	switch c.cfg.Storage.Kind {
	case config.Redis:
		c.storage = redis.NewRedis(c.ctx, c.cfg.Storage.Redis)
		err = c.storage.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cdiscount) initPool() error {
	pool, err := NewPool(c.ctx, c.cfg.Proxy)
	if err != nil {
		return err
	}
	c.ProxyPool = pool
	return nil
}

func (c *Cdiscount) initClient() error {
	opts := client.Option{
		Headless:                c.cfg.Client.Headless,
		BlinkSettings:           "imagesEnabled=false",
		UserAgent:               net.UserAgent,
		IgnoreCertificateErrors: true,
	}
	if c.cfg.Client.Headless {
		opts.WindowsHigh, opts.WindowsWith = 400, 400
	}

	cli := client.NewClient(c.ctx, opts)
	err := cli.NewTask().Do(c.ctx, "https://www.baidu.com")
	if err != nil {
		return err
	}
	c.cli = cli
	return nil
}

func (c *Cdiscount) initServe() {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	controller.RegistryTaskController(c, handler)
	controller.RegistryGoodController(c.store, handler)

	c.Serve = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.cfg.Web.Binding, c.cfg.Web.Port),
		Handler:      handler,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
}

func (c *Cdiscount) Start(stop <-chan struct{}) {

	go c.Serve.ListenAndServe()
	log.Infof("start at %v", c.Serve.Addr)
	go c.daemon()
	log.Infof("start daemon")

	<-stop

	c.Close()

	return
}

func (c *Cdiscount) Close() error {
	c.cancel()
	c.ProxyPool.Close()
	return log.Sync()
}