package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func (s *freeip) StartUrl() []string {
	return []string{
		"https://www.freeip.top/",
		"https://www.freeip.top/?page=2",
		"https://www.freeip.top/?page=3",
		"https://www.freeip.top/?page=4",
		"https://www.freeip.top/?page=5",
	}
}

func (s *freeip) GetReferer() string {
	return "https://www.freeip.top/"
}

type freeip struct {
	Spider
}

func (s *freeip) Cron() string {
	return "@every 2m"
}

func (s *freeip) Name() string {
	return "freeip"
}

func (s *freeip) Run() {
	getProxy(s)
}

func (s *freeip) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
