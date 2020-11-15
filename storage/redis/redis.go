package redis

import (
	"context"
	"path"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/lack-io/cirrus/config"
	"github.com/lack-io/cirrus/storage"
)

const (
	prefix = "/cirrus"

	// redis 检测 redis 状态的时间间隔
	pingInterval = time.Second * 5
)

type Redis struct {
	// ctx 控制 Redis 的停止
	ctx context.Context

	// cli redis 客户端
	cli *redis.Client

	// raw redis hash 表的名称，存储未被爬取过的 url
	raw string

	// cook redis hash 表名称，存储爬取过的 url
	cook string

	ready *atomic.Value
}

func NewRedis(ctx context.Context, cfg *config.StorageRedis) *Redis {
	cli := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           0,
		PoolSize:     cfg.Pools,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})

	rdb := &Redis{
		ctx:   ctx,
		cli:   cli,
		raw:   path.Join(prefix, "raw"),
		cook:  path.Join(prefix, "cook"),
		ready: &atomic.Value{},
	}

	rdb.ready.Store(false)

	return rdb
}

func (r *Redis) Init() error {
	c := r.cli.Ping(r.ctx)
	if c.Err() != nil {
		return c.Err()
	}

	r.cli.HSet(r.ctx, r.cook, "1")
	r.ready.Store(true)
	go r.ping()
	return nil
}

func (r *Redis) ping() {
	timer := time.NewTicker(pingInterval)
	for {
		select {
		case <-r.ctx.Done():
			timer.Stop()
			_ = r.cli.Close()
			return
		case <-timer.C:
			c := r.cli.Ping(r.ctx)
			if c.Err() != nil {
				r.ready.Store(false)
			}
		}
	}
}

func (r *Redis) Get() (*storage.URL, error) {
	if !r.ready.Load().(bool) {
		return nil, storage.ErrStorage
	}

	result := r.cli.SPop(r.ctx, r.raw)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return &storage.URL{Path: result.Val(), Storage: r}, nil
}

func (r *Redis) Push(url storage.URL) error {
	if !r.ready.Load().(bool) {
		return storage.ErrStorage
	}

	if r.cli.HExists(r.ctx, r.cook, url.Path).Val() {
		return storage.ErrOldURL
	}

	return r.cli.SAdd(r.ctx, r.raw, url.Path).Err()
}

func (r *Redis) Persist(url storage.URL) error {
	if !r.ready.Load().(bool) {
		return storage.ErrStorage
	}


	return r.cli.HSet(r.ctx, r.cook, url.Path, 1).Err()
}

func (r *Redis) Reset() {
	// 清除访问过的 URL
	r.cli.Del(r.ctx, r.cook)
	r.cli.Del(r.ctx, r.raw)
}
