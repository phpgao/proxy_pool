package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func (s *newproxy) StartUrl() []string {
	return []string{
		"http://newproxy.org.ru/page.php?page_id=1",
	}
}

func (s *newproxy) GetReferer() string {
	return "http://newproxy.org.ru"
}

type newproxy struct {
	Spider
}

func (s *newproxy) Cron() string {
	return "@every 10m"
}

func (s *newproxy) Name() string {
	return "newproxy"
}

func (s *newproxy) Run() {
	getProxy(s)
}

func (s *newproxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
