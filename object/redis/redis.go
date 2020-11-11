package redis

import (
	"context"
	"fmt"
	"path"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/lack-io/cirrus/object"
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

func NewRedis(ctx context.Context, addr, username, password string, size int) *Redis {
	cli := redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     username,
		Password:     password,
		DB:           0,
		PoolSize:     size,
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

func (r *Redis) GetURL() (*object.URL, error) {
	if !r.ready.Load().(bool) {
		return nil, object.ErrObject
	}

	num := r.cli.SCard(r.ctx, r.raw).Val()
	if num == 0 {
		return nil, object.ErrNoURL
	}

	c := r.cli.SPop(r.ctx, r.raw)
	if c.Err() != nil {
		return nil, fmt.Errorf("%w: %v", object.ErrGetURL, c.Err())
	}

	url := &object.URL{
		Path: c.Val(),
		Obj:  r,
	}

	r.cli.HSet(r.ctx, r.cook, c.Val(), 1)

	return url, nil
}

func (r *Redis) SetURL(url *object.URL) error {
	if !r.ready.Load().(bool) {
		return object.ErrObject
	}

	if r.cli.SIsMember(r.ctx, r.raw, url.Path).Val() {
		return object.ErrURLExists
	}

	if r.cli.HExists(r.ctx, r.cook, url.Path).Val() {
		return object.ErrOldURL
	}

	if err := r.cli.SAdd(r.ctx, r.raw, url.Path).Err(); err != nil {
		return fmt.Errorf("%w: %v", object.ErrSetURL, err)
	}

	return nil
}

func (r *Redis) Reset() {
	r.cli.Del(r.ctx, r.raw)
	r.cli.Del(r.ctx, r.cook)
}
