package pool

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var p *Pool

func TestNew(t *testing.T) {
	p = New(context.Background(), 2)
}

func TestPool_NewTask(t *testing.T) {
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("task%d", i)
		fn := func(name string) func() {
			return func() {
				now := time.Now()
				fmt.Printf("%s start at %v\n", name, now)
				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(10)
				time.Sleep(time.Second * time.Duration(n))
				fmt.Printf("%s done!!! -- speed:%v\n", name, time.Now().Sub(now))
			}
		}
		p.NewTask(fn(name))
	}
}
