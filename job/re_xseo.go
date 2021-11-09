package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func init() {
	spider := xseo{}
	if spider.Enabled() {
		register(spider)
	}
}

func (s xseo) StartUrl() []string {
	return []string{
		"http://xseo.in/freeproxy",
	}
}

func (s xseo) GetReferer() string {
	return "http://xseo.in"
}

type xseo struct {
	Spider
}

func (s xseo) Cron() string {
	return "@every 5m"
}

func (s xseo) Name() string {
	return "xseo"
}

func (s xseo) Run() {
	getProxy(s)
}

func (s xseo) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
