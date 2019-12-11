package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func (s *aliveProxy) StartUrl() []string {
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

func (s *aliveProxy) GetReferer() string {
	return "http://aliveproxy.com"
}

type aliveProxy struct {
	Spider
}

func (s *aliveProxy) Cron() string {
	return "@every 5m"
}

func (s *aliveProxy) Name() string {
	return "aliveproxy"
}

func (s *aliveProxy) Run() {
	getProxy(s)
}

func (s *aliveProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
