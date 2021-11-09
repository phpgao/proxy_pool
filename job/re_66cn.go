package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"regexp"
	"strings"
)

func init() {
	spider := cn66{}
	if spider.Enabled() {
		register(spider)
	}
}

func (s cn66) StartUrl() []string {
	return []string{
		"http://www.66ip.cn/mo.php?tqsl=2000",
		"http://www.66ip.cn/nmtq.php?getnum=300&isp=0&anonymoustype=0&start=&ports=&export=&ipaddress=&area=0&proxytype=2&api=66ip",
		"http://www.66ip.cn/mo.php?sxb=&tqsl=300&port=&export=&ktip=&sxa=&submit=%CC%E1++%C8%A1&textarea=",
	}
}

func (s cn66) GetReferer() string {
	return "http://www.66ip.cn/"
}

type cn66 struct {
	Spider
}

func (s cn66) Cron() string {
	return "@every 1m"
}

func (s cn66) Name() string {
	return "cn66"
}

func (s cn66) Run() {
	getProxy(s)
}

func (s cn66) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
