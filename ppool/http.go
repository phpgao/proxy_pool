package ppool

import (
	"errors"
	"github.com/fatih/pool"
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/util"
	"math/rand"
	"net"
	"time"
)

var (
	Http   pool.Pool
	Https  pool.Pool
	config = util.GetConfig()
	logger = util.GetLogger()
)

func init() {
	var err error
	Https, err = pool.NewChannelPool(0, 20, func() (conn net.Conn, err error) {

		storeEngine := db.GetDb()
		proxies, err := storeEngine.Get(map[string]string{
			"schema": "http",
		})
		if err != nil {
			return
		}
		l := len(proxies)
		if l == 0 {
			err = errors.New("no proxy")
			return
		}
		proxy := proxies[rand.Intn(l)]
		conn, err = net.DialTimeout("tcp", proxy.GetProxyUrl(), time.Duration(config.TcpTimeout)*time.Second)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		logger.Fatal("fatal create https pool")
	}
	Http, err = pool.NewChannelPool(0, 20, func() (conn net.Conn, err error) {
		storeEngine := db.GetDb()
		proxies, err := storeEngine.Get(map[string]string{
			"schema": "http",
		})
		if err != nil {
			return
		}
		l := len(proxies)
		if l == 0 {
			err = errors.New("no proxy")
			return
		}
		proxy := proxies[rand.Intn(l)]
		conn, err = net.DialTimeout("tcp", proxy.GetProxyUrl(), time.Duration(config.TcpTimeout)*time.Second)
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		logger.Fatal("fatal create http pool")
	}
}
