package validator

import (
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/queue"
	"github.com/phpgao/proxy_pool/util"
	"sync"
)

var (
	config      = util.GetConfig()
	logger      = util.GetLogger()
	storeEngine = db.GetDb()
	lockMap     = sync.Map{}
)

func NewValidator() {
	q := queue.GetNewChan()
	for {
		proxy := <-q
		go func(ip *model.HttpProxy) {
			key := ip.GetKey()
			if _, ok := lockMap.Load(key); ok {
				return
			}

			lockMap.Store(key, 1)
			defer func() {
				lockMap.Delete(key)
			}()

			if storeEngine.Exists(*proxy) {
				logger.WithField("proxy", proxy.GetKey()).WithField(
					"proxy", proxy.GetProxyUrl()).Debug("skip exists proxy")
				return
			}
			if !ip.SimpleTcpTest() {
				logger.WithField("proxy", proxy.GetProxyUrl()).Debug("failed tcp test")
				return
			}
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
}
