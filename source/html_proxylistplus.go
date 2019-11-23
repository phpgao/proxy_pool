package source

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type proxylistplus struct {
	Spider
}

func (s *proxylistplus) Cron() string {
	return "@every 5m"
}

func (s *proxylistplus) StartUrl() []string {
	return []string{
		"https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1",
	}
}

func (s *proxylistplus) GetReferer() string {
	return "https://list.proxylistplus.com/"
}

func (s *proxylistplus) Run() {
	getProxy(s)
}

func (s *proxylistplus) Name() string {
	return "proxylistplus"
}

func (s *proxylistplus) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	list := htmlquery.Find(doc, "//*[@id='page']/table[2]/tbody/tr[position()>2]")

	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[3]"))

		ip = strings.TrimSpace(ip)
		port = strings.TrimSpace(port)

		proxies = append(proxies, &model.HttpProxy{
			Ip:   ip,
			Port: port,
		})
	}

	return
}
