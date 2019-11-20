package spider

import (
	"fmt"
	"github.com/phpgao/proxy_pool/model"
	"net/url"
	"regexp"
	"strings"
)

var countryList = []string{
	"印度尼西亚",
	"美国",
	"泰国",
	"俄罗斯",
	"巴西",
	"乌克兰",
	"北美地区",
	"亚太地区",
	"印度",
	"IANA",
	"波兰",
	"英国",
	"哥伦比亚",
	"拉美地区",
	"欧洲和中东地区",
	"日本",
	"阿根廷",
	"孟加拉",
	"柬埔寨",
}

func (s *ip89) StartUrl() []string {
	var t []string
	for _, c := range countryList {
		t = append(t, fmt.Sprintf("http://www.89ip.cn/tqdl.html?api=1&num=100&port=&address%s&isp=", url.QueryEscape(c)))
	}
	return t
}

func (s *ip89) GetReferer() string {
	return "http://www.89ip.cn/"
}

type ip89 struct {
	Spider
}

func (s *ip89) Cron() string {
	return "@every 1m"
}

func (s *ip89) Name() string {
	return "ip89"
}

func (s *ip89) Run() {
	getProxy(s)
}

func (s *ip89) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
