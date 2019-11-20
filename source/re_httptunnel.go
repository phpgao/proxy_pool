package source

import (
	"github.com/phpgao/proxy_pool/model"
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
	return "@every 1m"
}

func (s *httptunnel) Name() string {
	return "httptunnel"
}

func (s *httptunnel) Run() {
	getProxy(s)
}

func (s *httptunnel) Parse(body string) (proxies []*model.HttpProxy, err error) {
	reg := regexp.MustCompile(`(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5]):\d{0,5}`)
	rs := reg.FindAllString(body, -1)

	for _, proxy := range rs {
		if strings.Contains(proxy, ":") {
			proxyInfo := strings.Split(proxy, ":")

			proxies = append(proxies, &model.HttpProxy{
				Ip:        proxyInfo[0],
				Port:      proxyInfo[1],
			})
		}
	}
	return
}
