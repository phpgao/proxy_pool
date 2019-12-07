package cache

import (
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"sync"
	"sync/atomic"
	"time"
)

type Cached struct {
	m       sync.RWMutex
	proxy   map[string][]model.HttpProxy
	expire  time.Time
	expired bool
}

var (
	logger = util.GetLogger("cache")
	Cache  = Cached{
		proxy: map[string][]model.HttpProxy{},
	}
	engine       = db.GetDb()
	once         Once
	cacheTimeout = time.Duration(util.ServerConf.ProxyCacheTimeOut)
)

func init() {
	Cache.proxy = getProxyMap()
	Cache.expire = time.Now().Add(cacheTimeout * time.Second)
	go func() {
		for {
			if time.Now().Sub(Cache.expire) >= 0 {
				atomic.StoreUint32(&once.done, 0)
			}
			time.Sleep(time.Second)
		}
	}()
}

func getProxyMap() map[string][]model.HttpProxy {
	m := map[string][]model.HttpProxy{
		"http":  nil,
		"https": nil,
	}
	var err error
	for k, _ := range m {
		var p []model.HttpProxy

		if k == "http" {
			p, err = engine.Get(map[string]string{
				"score": string(util.ServerConf.ScoreAtLeast),
			})
		} else {
			p, err = engine.Get(map[string]string{
				"score":  string(util.ServerConf.ScoreAtLeast),
				"schema": k,
			})
		}

		if err != nil {
			logger.WithField("type", k).WithError(err).Error("error get proxy")
			continue
		}
		m[k] = p
	}

	return m
}

func (c *Cached) Get() map[string][]model.HttpProxy {
	c.m.RLock()
	defer c.m.RUnlock()
	now := time.Now()
	if now.Sub(c.expire) >= 0 {
		once.Do(c.Update)
	}
	return Cache.proxy
}

func (c *Cached) Update() {
	logger.Debug("updating cache")
	c.proxy = getProxyMap()
	c.expire = time.Now().Add(cacheTimeout * time.Second)
}

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
