package source

import (
	"github.com/phpgao/proxy_pool/model"
	"regexp"
	"strings"
)

func (s *cn66) StartUrl() []string {
	return []string{
		"http://www.66ip.cn/mo.php?tqsl=1000",
	}
}

func (s *cn66) GetReferer() string {
	return "http://www.66ip.cn/"
}

type cn66 struct {
	Spider
}

func (s *cn66) Cron() string {
	return "@every 2m"
}

func (s *cn66) Name() string {
	return "cn66"
}

func (s *cn66) Run() {
	getProxy(s)
}

func (s *cn66) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
