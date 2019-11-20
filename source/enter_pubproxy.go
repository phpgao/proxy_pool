package source

import (
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type pubProxy struct {
	Spider
}

func (s *pubProxy) StartUrl() []string {
	return []string{
		"http://pubproxy.com/api/proxy?limit=5&format=txt&type=http&level=anonymous&last_check=60&no_country=CN",
		"http://pubproxy.com/api/proxy?limit=5&format=txt&type=http&level=anonymous&last_check=60&country=CN",
	}
}

func (s *pubProxy) GetReferer() string {
	return "http://pubproxy.com/"
}

func (s *pubProxy) Cron() string {
	return "@every 1m"
}

func (s *pubProxy) Name() string {
	return "pub"
}

func (s *pubProxy) Run() {
	getProxy(s)
}

func (s *pubProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	proxyStrings := strings.Split(body, "\n")
	for _, proxy := range proxyStrings {
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
