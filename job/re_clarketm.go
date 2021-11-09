package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func init() {
	spider := clarketm{}
	if spider.Enabled() {
		register(spider)
	}
}

func (s clarketm) StartUrl() []string {
	return []string{
		"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list.txt",
	}
}

func (s clarketm) GetReferer() string {
	return "http://github.com/"
}

type clarketm struct {
	Spider
}

func (s clarketm) Cron() string {
	return "@every 5m"
}

func (s clarketm) Name() string {
	return "clarketm"
}

func (s clarketm) Run() {
	getProxy(s)
}

func (s clarketm) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
