package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func (s *blackHat) StartUrl() []string {
	return []string{
		"http://www.blackhat.be/cpt/proxy.lst",
	}
}

func (s *blackHat) GetReferer() string {
	return "http://www.blackhat.be/cpt/proxy.lst"
}

type blackHat struct {
	Spider
}

func (s *blackHat) Cron() string {
	return "@every 1h"
}

func (s *blackHat) Name() string {
	return "blackHat"
}

func (s *blackHat) Run() {
	getProxy(s)
}

func (s *blackHat) Parse(body string) (proxies []*model.HttpProxy, err error) {
	reg := regexp.MustCompile(util.RegProxyWithoutColon)
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
