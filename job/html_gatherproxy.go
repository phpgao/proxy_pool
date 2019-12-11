package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type gatherproxy struct {
	Spider
}

func (s *gatherproxy) StartUrl() []string {
	return []string{
		"http://www.gatherproxy.com/proxylist/country/?c=China",
	}
}

func (s *gatherproxy) Cron() string {
	return "@every 30m"
}

func (s *gatherproxy) GetReferer() string {
	return "http://free-proxy.gatherproxy/zh/proxylist/country/CN/all/ping/all"
}

func (s *gatherproxy) Run() {
	getProxy(s)
}

func (s *gatherproxy) Name() string {
	return "gatherproxy"
}

func (s *gatherproxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
