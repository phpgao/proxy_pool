package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func (s *dogdev) StartUrl() []string {
	return []string{
		"http://dogdev.net/Proxy/all",
	}
}

func (s *dogdev) GetReferer() string {
	return "http://dogdev.net/"
}

type dogdev struct {
	Spider
}

func (s *dogdev) Cron() string {
	return "@every 1h"
}

func (s *dogdev) Name() string {
	return "dogdev"
}

func (s *dogdev) Run() {
	getProxy(s)
}

func (s *dogdev) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
