package validator

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/queue"
)

func OldValidator() {
	q := queue.GetOldChan()
	for {
		ip := <-q
		go func(ip model.HttpProxy) {
			key := ip.GetKey()
			if _, ok := lockMap.Load(key); ok {
				return
			}

			lockMap.Store(key, 1)
			defer func() {
				lockMap.Delete(key)
			}()
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
}
