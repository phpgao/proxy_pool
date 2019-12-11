package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type ip3366 struct {
	Spider
}

func (s *ip3366) StartUrl() []string {
	return []string{
		"http://www.ip3366.net/free/?stype=1",
		"http://www.ip3366.net/free/?stype=1&page=2",
		"http://www.ip3366.net/free/?stype=2",
		"http://www.ip3366.net/free/?stype=2&page=2",
		"http://proxy.ip3366.net/free/?action=china&page=1",
		"http://proxy.ip3366.net/free/?action=china&page=2",
		"http://proxy.ip3366.net/free/?action=china&page=3",
	}
}

func (s *ip3366) Cron() string {
	return "@every 30m"
}

func (s *ip3366) GetReferer() string {
	return "http://www.ip3366.net"
}

func (s *ip3366) Run() {
	getProxy(s)
}

func (s *ip3366) Name() string {
	return "Kuai"
}

func (s *ip3366) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	list := htmlquery.Find(doc, "//table/tbody/tr[position()>1]")
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
