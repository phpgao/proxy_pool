package source

import (
	"github.com/phpgao/proxy_pool/model"
	"regexp"
	"strings"
)

func (s *newproxy) StartUrl() []string {
	return []string{
		"http://newproxy.org.ru/page.php?page_id=1",
	}
}

func (s *newproxy) GetReferer() string {
	return "http://newproxy.org.ru"
}

type newproxy struct {
	Spider
}

func (s *newproxy) Cron() string {
	return "@every 10m"
}

func (s *newproxy) Name() string {
	return "newproxy"
}

func (s *newproxy) Run() {
	getProxy(s)
}

func (s *newproxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
