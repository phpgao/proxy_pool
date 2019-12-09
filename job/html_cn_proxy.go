package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type cn_proxy struct {
	Spider
}

func (s *cn_proxy) StartUrl() []string {
	return []string{
		"http://cn-proxy.com/",
	}
}

func (s *cn_proxy) Cron() string {
	return "@every 30m"
}

func (s *cn_proxy) GetReferer() string {
	return "http://free-proxy.cn_proxy/zh/proxylist/country/CN/all/ping/all"
}

func (s *cn_proxy) Run() {
	getProxy(s)
}

func (s *cn_proxy) Name() string {
	return "cn_proxy"
}

func (s *cn_proxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	list := htmlquery.Find(doc, "//*[@id='list']/table/tbody/tr[position()>1]")
	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))

		ip = strings.TrimSpace(ip)
		port = strings.TrimSpace(port)

		proxies = append(proxies, &model.HttpProxy{
			Ip:   ip,
			Port: port,
		})
	}
	return
}
