package source

import (
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

func (s *ab57) StartUrl() []string {
	return []string{
		"http://ab57.ru/downloads/proxyold.txt",
	}
}

func (s *ab57) GetReferer() string {
	return "http://ab57.ru/"
}

type ab57 struct {
	Spider
}

func (s *ab57) Cron() string {
	return "@every 2m"
}

func (s *ab57) Name() string {
	return "ab57"
}

func (s *ab57) Run() {
	getProxy(s)
}

func (s *ab57) Parse(body string) (proxies []*model.HttpProxy, err error) {
	proxyString := strings.Split(body, "\r\n")
	for _, proxy := range proxyString {
		if strings.Contains(proxy, ":") {
			proxyInfo := strings.Split(proxy, ":")
			proxies = append(proxies, &model.HttpProxy{
				Ip:        proxyInfo[0],
				Port:      proxyInfo[1],
				Anonymous: 0,
			})
		}
	}
	return
}
