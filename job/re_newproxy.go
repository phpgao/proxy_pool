package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func init() {
	spider := newProxy{}
	if spider.Enabled() {
		register(spider)
	}
}

func (s newProxy) StartUrl() []string {
	return []string{
		"http://newproxy.org.ru/page.php?page_id=1",
	}
}

func (s newProxy) GetReferer() string {
	return "http://newproxy.org.ru"
}

type newProxy struct {
	Spider
}

func (s newProxy) Cron() string {
	return "@every 10m"
}

func (s newProxy) Name() string {
	return "newproxy"
}

func (s newProxy) Run() {
	getProxy(s)
}

func (s newProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
