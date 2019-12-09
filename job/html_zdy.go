package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type zdy struct {
	Spider
}

func (s *zdy) StartUrl() []string {
	return []string{
		"https://www.zdaye.com/dayProxy/ip/317476.html",
	}
}

func (s *zdy) Cron() string {
	return "@every 30m"
}

func (s *zdy) GetReferer() string {
	return "http://free-proxy.zdy/zh/proxylist/country/CN/all/ping/all"
}

func (s *zdy) Run() {
	getProxy(s)
}

func (s *zdy) Name() string {
	return "zdy"
}

func (s *zdy) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
