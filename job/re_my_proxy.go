package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func init() {
	spider := myProxy{}
	if spider.Enabled() {
		register(spider)
	}
}

func (s myProxy) StartUrl() []string {
	return []string{
		"https://www.my-proxy.com/free-proxy-list.html",
		"https://www.my-proxy.com/free-proxy-list-2.html",
	}
}

func (s myProxy) GetReferer() string {
	return "https://www.my-proxy.com/"
}

type myProxy struct {
	Spider
}

func (s myProxy) Cron() string {
	return "@every 10m"
}

func (s myProxy) Name() string {
	return "my_proxy"
}

func (s myProxy) Run() {
	getProxy(s)
}

func (s myProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
