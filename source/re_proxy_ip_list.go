package source

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func (s *proxyIpList) StartUrl() []string {
	return []string{
		"http://proxy-ip-list.com/download/free-proxy-list",
	}
}

func (s *proxyIpList) GetReferer() string {
	return "https://www.my-proxy.com/"
}

type proxyIpList struct {
	Spider
}

func (s *proxyIpList) Cron() string {
	return "@every 10m"
}

func (s *proxyIpList) Name() string {
	return "proxy-ip-list"
}

func (s *proxyIpList) Run() {
	getProxy(s)
}

func (s *proxyIpList) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
