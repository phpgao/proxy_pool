package validator

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/queue"
	"sync"
)

func OldValidator() {
	q := queue.GetOldChan()
	var wg sync.WaitGroup

	for i := 0; i < config.OldQueue; i++ {
		wg.Add(1)
		go func() {
			for {
				proxy := <-q
				func(p model.HttpProxy) {
					key := p.GetKey()
					if _, ok := lockMap.Load(key); ok {
						return
					}

					lockMap.Store(key, 1)
					defer func() {
						lockMap.Delete(key)
					}()
					if storeEngine.Exists(p) {
						var score = -30
						flag := p.SimpleTcpTest(config.GetTcpTestTimeOut())
						if flag {
							score = 10
						}
						err := storeEngine.AddScore(p, score)
						if err != nil {
							logger.WithError(err).WithField(
								"proxy", p.GetProxyUrl()).Error("error setting score")
						}
					}
				}(*proxy)
			}

		}()

	}
	wg.Wait()
}
