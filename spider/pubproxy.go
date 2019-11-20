package spider

import (
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type PubProxy struct {
	Spider
}

func (s *PubProxy) StartUrl() []string {
	return []string{
		"http://pubproxy.com/api/proxy?limit=5&format=txt&type=http&level=anonymous&last_check=60&no_country=CN",
		"http://pubproxy.com/api/proxy?limit=5&format=txt&type=http&level=anonymous&last_check=60&country=CN",
	}
}

func (s *PubProxy) RandomDelay() bool {
	return true
}

func (s *PubProxy) GetReferer() string {
	return "http://pubproxy.com/"
}

func (s *PubProxy) Cron() string {
	return "@every 1m"
}

func (s *PubProxy) Name() string {
	return "pub"
}

func (s *PubProxy) Run() {
	getProxy(s)
}

func (s *PubProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	proxyStrings := strings.Split(body, "\n")
	for _, proxy := range proxyStrings {
		if strings.Contains(proxy, ":") {
			proxyInfo := strings.Split(proxy, ":")
			proxies = append(proxies, &model.HttpProxy{
				Ip:        proxyInfo[0],
				Port:      proxyInfo[1],
				Schema:    "http",
				Score:     config.Score,
				Latency:   0,
				From:      s.Name(),
				Anonymous: 0,
			})
		}
	}
	return
}