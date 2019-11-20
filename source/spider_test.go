package source

import (
	"fmt"
	"github.com/apex/log"
	"github.com/phpgao/proxy_pool/model"
	"testing"
	"time"
)

func init() {
	logger.Level = log.DebugLevel
}

func TestGetSpiders(t *testing.T) {
	//for _, c := range ListOfSpider {
	//	testSpiderFetch(c)
	//}
	testSpider := &httptunnel{}
	testSpiderFetch(testSpider)
}

func testSpiderFetch(c Crawler) {
	newProxyChan := make(chan *model.HttpProxy, 100)
	c.SetProxyChan(newProxyChan)
	c.Run()
	timeout := time.After(30 * time.Second)
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
