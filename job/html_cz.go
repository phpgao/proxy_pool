package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

func init() {
	spider := cz{}
	if spider.Enabled() {
		register(spider)
	}
}

type cz struct {
	Spider
}

func (s cz) StartUrl() []string {
	return []string{
		"http://free-proxy.cz/zh/proxylist/country/CN/all/ping/all",
		"http://free-proxy.cz/zh/proxylist/country/CN/all/ping/all/2",
		"http://free-proxy.cz/zh/proxylist/country/CN/all/ping/all/3",
		"http://free-proxy.cz/zh/proxylist/country/CN/all/ping/all/4",
		"http://free-proxy.cz/zh/proxylist/country/CN/all/ping/all/5",
	}
}

func (s cz) Cron() string {
	return "@every 30m"
}

func (s cz) GetReferer() string {
	return "http://free-proxy.cz/zh/proxylist/country/CN/all/ping/all"
}

func (s cz) Run() {
	getProxy(s)
}

func (s cz) Name() string {
	return "cz"
}

func (s cz) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
