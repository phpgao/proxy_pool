package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

func init() {
	spider := kxdaili{}
	if spider.Enabled() {
		register(spider)
	}
}

type kxdaili struct {
	Spider
}

func (s kxdaili) StartUrl() []string {
	return []string{
		"http://www.kxdaili.com/dailiip.html",
		"http://www.kxdaili.com/dailiip/2/1.html",
	}
}

func (s kxdaili) Cron() string {
	return "@every 30m"
}

func (s kxdaili) GetReferer() string {
	return "http://free-proxy.kxdaili/zh/proxylist/country/CN/all/ping/all"
}

func (s kxdaili) Run() {
	getProxy(s)
}

func (s kxdaili) Name() string {
	return "kxdaili"
}

func (s kxdaili) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
