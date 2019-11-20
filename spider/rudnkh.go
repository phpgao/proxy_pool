package spider

import (
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

func (s *Rudnkh) StartUrl() []string {
	return []string{
		"https://proxy.rudnkh.me/txt",
		"https://raw.githubusercontent.com/a2u/free-proxy-list/master/free-proxy-list.txt",
	}
}

func (s *Rudnkh) GetReferer() string {
	return "https://proxy.rudnkh.me/"
}

type Rudnkh struct {
	Spider
}

func (s *Rudnkh) Cron() string {
	return "@every 1m"
}

func (s *Rudnkh) Name() string {
	return "rudnkh"
}

func (s *Rudnkh) Run() {
	getProxy(s)
}

func (s *Rudnkh) Parse(body string) (proxies []*model.HttpProxy, err error) {

	proxyString := strings.Split(body, "\n")
	for _, proxy := range proxyString {
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
