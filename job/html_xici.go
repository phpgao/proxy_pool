package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

func init() {
	spider := xici{}
	if spider.Enabled() {
		register(spider)
	}
}

type xici struct {
	Spider
}

func (s xici) Cron() string {
	return "@every 2m"
}

func (s xici) GetReferer() string {
	return "https://www.xicidaili.com/"
}

func (s xici) StartUrl() []string {
	return []string{
		"http://www.xicidaili.com/nn",
		"http://www.xicidaili.com/wn",
	}
}

func (s xici) Run() {
	getProxy(s)
}

func (s xici) Name() string {
	return "xici"
}

func (s xici) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	list := htmlquery.Find(doc, "//*[@id='ip_list']/tbody/tr[position()>1]")
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
