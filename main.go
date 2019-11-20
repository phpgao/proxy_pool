package main

import (
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/schedule"
	"github.com/phpgao/proxy_pool/server"
	"github.com/phpgao/proxy_pool/util"
)

var (
	config = util.GetConfig()
	logger = util.GetLogger()
)

// Check new proxy with tcp
var SimpleCheckProxyChan chan *model.HttpProxy

// Check new proxy with test
var newProxyChan chan *model.HttpProxy

// Present proxy in and out
var oldProxyChan chan *model.HttpProxy

var storeEngine db.Store

func init() {
	storeEngine = db.GetDb()
	newProxyChan = make(chan *model.HttpProxy, config.Concurrence)
	SimpleCheckProxyChan = make(chan *model.HttpProxy, config.Concurrence)
	oldProxyChan = make(chan *model.HttpProxy, config.Concurrence)
}

func main() {

	// fetch && feed validation queue
	if config.Manager {
		s := schedule.NewScheduler(SimpleCheckProxyChan, oldProxyChan, storeEngine)
		go s.Run()
	}

	go func() {
		for {
			ip := <-SimpleCheckProxyChan
			go func(ip *model.HttpProxy) {
				//logger.WithField("proxy", ip.GetProxyUrl()).Debug("simple test")
				if ip.SimpleTcpTest() {
					newProxyChan <- ip
				}
			}(ip)
		}
	}()

	// test new proxy
	go func() {
		for {
			proxy := <-newProxyChan
			if storeEngine.Exists(*proxy) {
				logger.WithField("proxy", proxy.GetKey()).WithField(
					"proxy", proxy.GetProxyUrl()).Debug("skip exists proxy")
				continue
			}
			go func(ip *model.HttpProxy) {
				// http test
				err := ip.TestProxy(false)
				if err != nil {
					logger.WithError(err).WithField(
						"proxy", ip.GetProxyUrl()).Debug("error when test http proxy")
				} else {

					logger.WithField("proxy", ip.GetProxyUrl()).WithField(
						"key", ip.GetKey()).Info("valid proxy")
					// https test
					err := ip.TestProxy(true)
					if err != nil {
						logger.WithError(err).WithField(
							"proxy", ip.GetProxyUrl()).Debug("error when test https proxy")
					}

					storeEngine.Add(*ip)
				}
			}(proxy)
		}
	}()
	// test old proxy
	go func() {
		for {
			ip := <-oldProxyChan
			//logger.WithField("proxy", ip.GetProxyUrl()).WithField(
			//	"key", ip.GetKey()).Debug("oldProxyChan")
			go func(ip model.HttpProxy) {
				if storeEngine.Exists(ip) {
					flag := ip.SimpleTcpTest()
					if flag {
						err := ip.TestProxy(false)
						if err == nil {
							err = storeEngine.AddScore(ip, 10)
							if err != nil {
								logger.WithError(err).WithField("proxy", ip.GetProxyUrl()).Error("error add score")
							}
							return
						}
					}
					err := storeEngine.AddScore(ip, -20)
					if err != nil {
						logger.WithError(err).WithField(
							"proxy", ip.GetProxyUrl()).Error("error when add score")
					}
				}
			}(*ip)
		}
	}()
	// run http server
	server.Serve()
}
