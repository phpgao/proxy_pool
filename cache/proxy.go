package cache

import (
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"sync"
	"time"
)

type Cached struct {
	m       sync.RWMutex
	proxy   []model.HttpProxy
	c       chan int
	expire  time.Time
	expired bool
}

var (
	Cache = Cached{
		c: make(chan int, 1000),
	}
	engine = db.GetDb()
)

func init() {
	Cache.proxy = engine.GetAll()
	Cache.expire = time.Now().Add(5 * time.Second)
	go Cache.Sync()
}

func (c *Cached) Get() []model.HttpProxy {
	c.m.RLock()
	defer c.m.RUnlock()
	now := time.Now()
	if now.Sub(c.expire) >= 0 {
		c.expired = true
		c.c <- 1
	}
	return Cache.proxy
}

func (c *Cached) Update() {
	c.m.Lock()
	defer c.m.Unlock()
	c.proxy = engine.GetAll()
	c.expire = time.Now().Add(5 * time.Second)
}

func (c *Cached) Sync() {
	for {
		select {
		case <-c.c:
			if c.expired {
				c.Update()
				c.expired = false
			}
		}

	}
}
