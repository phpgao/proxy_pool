package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type data5u struct {
	Spider
}

func (s *data5u) StartUrl() []string {
	return []string{
		"http://www.data5u.com/",
	}
}

func (s *data5u) Cron() string {
	return "@every 30m"
}

func (s *data5u) GetReferer() string {
	return "http://free-proxy.data5u/zh/proxylist/country/CN/all/ping/all"
}

func (s *data5u) Run() {
	getProxy(s)
}

func (s *data5u) Name() string {
	return "data5u"
}

func (s *data5u) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
