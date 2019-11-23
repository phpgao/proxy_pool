package source

import (
	"github.com/phpgao/proxy_pool/model"
	"regexp"
	"strings"
)

func (s *aliveproxy) StartUrl() []string {
	return []string{
		"http://aliveproxy.com/proxy-list-port-8080/",
		"http://aliveproxy.com/high-anonymity-proxy-list/",
		"http://aliveproxy.com/anonymous-proxy-list/",
		"http://aliveproxy.com/transparent-proxy-list/",
		"http://aliveproxy.com/proxy-list-port-80/",
		"http://aliveproxy.com/proxy-list-port-81/",
		"http://aliveproxy.com/proxy-list-port-3128/",
		"http://aliveproxy.com/proxy-list-port-8000/",
		"http://aliveproxy.com/proxy-list-port-8080/",
		"http://aliveproxy.com/proxy-list-port-8080/",
		"http://aliveproxy.com/us-proxy-list/",
		"http://aliveproxy.com/ru-proxy-list/",
		"http://aliveproxy.com/jp-proxy-list/",
		"http://aliveproxy.com/ca-proxy-list/",
	}
}

func (s *aliveproxy) GetReferer() string {
	return "http://aliveproxy.com"
}

type aliveproxy struct {
	Spider
}

func (s *aliveproxy) Cron() string {
	return "@every 5m"
}

func (s *aliveproxy) Name() string {
	return "aliveproxy"
}

func (s *aliveproxy) Run() {
	getProxy(s)
}

func (s *aliveproxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	reg := regexp.MustCompile(`(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5]):\d{0,5}`)
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
