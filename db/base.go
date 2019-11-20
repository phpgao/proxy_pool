package db

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
)

var (
	config = util.GetConfig()
	logger = util.GetLogger()
	db     Store
)

type Store interface {
	Init() error
	GetAll() []model.HttpProxy
	Get(map[string]string) ([]model.HttpProxy, error)
	Exists(model.HttpProxy) bool
	Add(model.HttpProxy) bool
	Remove(model.HttpProxy) (bool, error)
	Random() (model.HttpProxy, error)
	Len() int
	Test() bool
	AddScore(key model.HttpProxy, score int) error
}

func GetDb() Store {
	if db == nil {
		db = &redisDB{
			client: redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", config.Redis.IP, config.Redis.Port),
				Password: config.Redis.Auth, // no password set
				DB:       config.Redis.Db,   // use default DB
				//PoolSize:       10,   // use default DB
			}),
			PrefixKey: config.Redis.PrefixKey,
			KeyExpire: config.Expire,
		}
		if !db.Test() {
			panic("db test error")
		}
	}

	return db
}
