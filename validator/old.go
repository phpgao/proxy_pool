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
				var score = -30
				flag := ip.SimpleTcpTest(config.GetTcpTestTimeOut())
				if flag {
					score = 10
				}
				err := storeEngine.AddScore(ip, score)
				if err != nil {
					logger.WithError(err).WithField(
						"proxy", ip.GetProxyUrl()).Error("error setting score")
				}
			}
		}(*ip)
	}
}
