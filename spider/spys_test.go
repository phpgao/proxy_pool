package spider

import (
	"fmt"
	"github.com/phpgao/proxy_pool/model"
	"testing"
	"time"
)

func TestSpys_Run(t *testing.T) {
	newProxyChan := make(chan *model.HttpProxy, 100)
	spider := Spys{Spider{ch: newProxyChan}}
	spider.Run()
	timeout := time.After(40 * time.Second)
	for {
		select {
		case proxy := <-newProxyChan:
			fmt.Println(proxy)
			go func(p *model.HttpProxy) {
				flag := p.SimpleTcpTest()
				fmt.Println(flag)
			}(proxy)
		case <-timeout:
			fmt.Println("There's no more time to this. Exiting!")
			return
		}
	}
}
