package db

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type redisDB struct {
	PrefixKey string
	client    *redis.Client
	KeyExpire int
	lock      sync.RWMutex
}

func (r *redisDB) Test() bool {
	pong, err := r.client.Ping().Result()
	if err != nil {
		logger.WithError(err).Error("error test redis")
		return false
	}

	if pong != "PONG" {
		logger.WithField("pong", pong).Error("error connecting to redis")
		return false
	}
	return true
}

func (r *redisDB) Init() error {
	return nil
}

func (r *redisDB) GetProxyKey(proxy model.HttpProxy) string {
	return strings.Join([]string{
		r.PrefixKey,
		"list",
		proxy.GetKey(),
	}, ":")
}

func (r *redisDB) Add(proxy model.HttpProxy) bool {
	key := strings.Join([]string{
		r.PrefixKey,
		"list",
		proxy.GetKey(),
	}, ":")
	if !r.KeyExists(key) {
		err := r.client.HMSet(key, proxy.GetProxyMap()).Err()
		if err != nil {
			logger.WithError(err).Error("error add proxy")
			return false
		}
	} else {
		err := r.AddScore(proxy, 10)
		if err != nil {
			logger.WithError(err).Error("error add Score")
			return false
		}
	}
	// add ttl
	err := r.ExpireDefault(proxy)
	if err != nil {
		logger.WithError(err).Error("error setting expire")
		return false
	}

	return true
}

func (r *redisDB) UpdateSchema(proxy model.HttpProxy) (err error) {
	key := strings.Join([]string{
		r.PrefixKey,
		"list",
		proxy.GetKey(),
	}, ":")
	if !r.KeyExists(key) {
		return errors.New("proxy not exists")
	}

	err = r.client.HSet(key, "Schema", proxy.Schema).Err()
	if err != nil {
		return err
	}

	return
}

func (r *redisDB) Expire(key string, expiration time.Duration) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.client.Expire(key, expiration).Err()
}
func (r *redisDB) SetScore(p model.HttpProxy) error {
	key := r.GetProxyKey(p)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.client.HSet(key, "Score", p.Score).Err()
}
func (r *redisDB) GetScore(p model.HttpProxy) int {
	key := r.GetProxyKey(p)
	v, err := r.client.HGet(key, "Score").Int()
	if err != nil {
		logger.WithField("key", key).WithError(err).Error("error get score")
		return 0
	}
	return v
}

func (r *redisDB) ExpireDefault(p model.HttpProxy) error {
	key := r.GetProxyKey(p)
	if r.KeyExpire <= 0 {
		return nil
	}
	return r.client.Expire(key, time.Duration(r.KeyExpire)*time.Second).Err()
}

func (r *redisDB) AddScore(p model.HttpProxy, score int) (err error) {

	v := r.GetScore(p)
	rs := v + score

	if rs <= 0 {
		err = r.Remove(p)
		if err != nil {
			return
		}
		return
	}
	if rs >= 100 {
		rs = 100
		err = r.ExpireDefault(p)
		if err != nil {
			return
		}
	}

	p.Score = rs

	return r.SetScore(p)
}

func (r *redisDB) KeyExists(key string) bool {
	val := r.client.Exists(key).Val()
	if val == 1 {
		return true
	}
	return false
}

func (r *redisDB) Exists(proxy model.HttpProxy) bool {
	key := r.GetProxyKey(proxy)
	return r.KeyExists(key)
}

func (r *redisDB) GetAll() (proxies []model.HttpProxy) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	keyPattern := strings.Join([]string{
		r.PrefixKey,
		"list",
		"*",
	}, ":")
	keys, _ := r.client.Keys(keyPattern).Result()
	for _, key := range keys {
		proxy := r.client.HGetAll(key).Val()
		//logger.WithField("proxy", proxy).Info("get all proxy")
		newProxy, err := model.Make(proxy)
		if err != nil {
			logger.WithField("key", key).WithError(err).Error("error create proxy from map")
			r.client.Del(key)
			continue
		}
		proxies = append(proxies, newProxy)
	}
	return proxies
}

func (r *redisDB) Get(options map[string]string) (proxies []model.HttpProxy, err error) {
	all := r.GetAll()
	filters, err := model.GetNewFilter(options)
	var limit int
	if l, ok := options["limit"]; ok {
		limit, err = strconv.Atoi(l)
		if err != nil {
			return
		}
	}

	if limit == 0 {
		limit = util.ServerConf.Limit
	}

	if err != nil {
		return
	}

	if len(filters) > 0 {
		for _, p := range all {
			if Match(filters, p) {
				proxies = append(proxies, p)
			}
			if len(proxies) > limit {
				return
			}
		}
	} else {
		if len(all) <= limit {
			proxies = all
		} else {
			proxies = all[:limit]
		}

	}
	return
}

func Match(filters []func(model.HttpProxy) bool, p model.HttpProxy) bool {
	for _, fc := range filters {
		if !fc(p) {
			return false
		}
	}

	return true
}

func (r *redisDB) Remove(proxy model.HttpProxy) (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	key := r.GetProxyKey(proxy)
	if r.KeyExists(key) {
		err = r.client.Del(key).Err()
		if err != nil {
			logger.WithError(err).WithField("key", key).Error("error deleting key")
			return err
		}
	}
	return nil
}

func (r *redisDB) RemoveAll(proxies []model.HttpProxy) (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	pipe := r.client.Pipeline()
	for _, proxy := range proxies {
		key := r.GetProxyKey(proxy)
		if r.KeyExists(key) {
			pipe.Del(key)
		}
	}
	_, err = pipe.Exec()
	if err != nil {
		return
	}
	return
}

func (r *redisDB) Random() (p model.HttpProxy, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	keyPattern := strings.Join([]string{
		r.PrefixKey,
		"list",
		"*",
	}, ":")
	keys, err := r.client.Keys(keyPattern).Result()
	if err != nil {
		logger.WithError(err).WithField("keys", keys).Error("error get keys from redis")
		return
	}
	if len(keys) == 0 {
		err = errors.New("no proxy")
		return
	}
	//rand.Seed(time.Now().Unix())
	key := keys[rand.Intn(len(keys))]
	proxy := r.client.HGetAll(key).Val()
	//logger.WithField("proxy", proxy).Info("get all proxy")
	newProxy, err := model.Make(proxy)
	if err != nil {
		return
	}
	return newProxy, nil
}

func (r *redisDB) Len() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	keyPattern := strings.Join([]string{
		r.PrefixKey,
		"list",
		"*",
	}, ":")
	keys, _ := r.client.Keys(keyPattern).Result()
	return len(keys)
}
