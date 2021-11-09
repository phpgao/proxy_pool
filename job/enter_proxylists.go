package job

import (
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

func init() {
	spider := proxyListsLine{}
	if spider.Enabled() {
		register(spider)
	}
}

func (s proxyListsLine) StartUrl() []string {
	return []string{
		"http://www.proxylists.net/http.txt",
		"http://www.proxylists.net/http_highanon.txt",
	}
}

func (s proxyListsLine) GetReferer() string {
	return "https://www.proxylists.net/"
}

type proxyListsLine struct {
	Spider
}

func (s proxyListsLine) Cron() string {
	return "@every 10m"
}

func (s proxyListsLine) Name() string {
	return "proxy-lists"
}

func (s proxyListsLine) Run() {
	getProxy(s)
}

func (s proxyListsLine) Parse(body string) (proxies []*model.HttpProxy, err error) {

	proxyString := strings.Split(body, "\n")
	for _, proxy := range proxyString {
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
