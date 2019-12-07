package source

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func (s *httptunnel) StartUrl() []string {
	return []string{
		"http://www.httptunnel.ge/ProxyListForFree.aspx",
	}
}

func (s *httptunnel) GetReferer() string {
	return "http://www.httptunnel.ge/"
}

type httptunnel struct {
	Spider
}

func (s *httptunnel) Cron() string {
	return "@every 5m"
}

func (s *httptunnel) Name() string {
	return "httptunnel"
}

func (s *httptunnel) Run() {
	getProxy(s)
}

func (s *httptunnel) Parse(body string) (proxies []*model.HttpProxy, err error) {
	reg := regexp.MustCompile(util.RegProxy)
	rs := reg.FindAllString(body, -1)

	for _, proxy := range rs {
		if strings.Contains(proxy, ":") {
			proxyInfo := strings.Split(proxy, ":")

			proxies = append(proxies, &model.HttpProxy{
				Ip:   proxyInfo[0],
				Port: proxyInfo[1],
			})
		}
	}
	return
}
