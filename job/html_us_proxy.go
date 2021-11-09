package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

func init() {
	spider := usProxy{}
	if spider.Enabled() {
		register(spider)
	}
}

type usProxy struct {
	Spider
}

func (s usProxy) Cron() string {
	return "@every 2m"
}

func (s usProxy) GetReferer() string {
	return "https://www.us-proxy.org/"
}

func (s usProxy) StartUrl() []string {
	return []string{
		"https://www.us-proxy.org/",
	}
}

func (s usProxy) Run() {
	getProxy(s)
}

func (s usProxy) Name() string {
	return "usProxy"
}

func (s usProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	list := htmlquery.Find(doc, "//*[@id='proxylisttable']/tbody/tr")
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
