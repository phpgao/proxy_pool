package source

import (
	"github.com/phpgao/proxy_pool/model"
	"regexp"
	"strings"
)

func (s *clarketm) StartUrl() []string {
	return []string{
		"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list.txt",
	}
}

func (s *clarketm) GetReferer() string {
	return "http://github.com/"
}

type clarketm struct {
	Spider
}

func (s *clarketm) Cron() string {
	return "@every 1m"
}

func (s *clarketm) Name() string {
	return "clarketm"
}

func (s *clarketm) Run() {
	getProxy(s)
}

func (s *clarketm) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
