package spider

import (
	"github.com/phpgao/proxy_pool/model"
	"regexp"
	"strings"
)

func (s *jiangxianli) StartUrl() []string {
	return []string{
		"http://ip.jiangxianli.com/",
	}
}

func (s *jiangxianli) GetReferer() string {
	return "http://ip.jiangxianli.com//"
}

type jiangxianli struct {
	Spider
}

func (s *jiangxianli) Cron() string {
	return "@every 30s"
}

func (s *jiangxianli) Name() string {
	return "jiangxianli"
}

func (s *jiangxianli) Run() {
	getProxy(s)
}

func (s *jiangxianli) Parse(body string) (proxies []*model.HttpProxy, err error) {
	reg := regexp.MustCompile(`(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5]):\d{0,5}`)
	rs := reg.FindAllString(body, -1)

	for _, proxy := range rs {
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
