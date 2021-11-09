package job

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

func init() {
	spider := nimadaili{}
	if spider.Enabled() {
		register(spider)
	}
}

type nimadaili struct {
	Spider
}

func (s nimadaili) StartUrl() []string {
	var u []string
	for _, d := range []string{"gaoni", "http", "https", "putong"} {
		for i := 1; i < 5; i++ {
			u = append(u, fmt.Sprintf("http://www.nimadaili.com/%s/%d/", d, i))
		}
	}
	return u
}

func (s nimadaili) Cron() string {
	return "@every 2m"
}

func (s nimadaili) GetReferer() string {
	return "http://www.nimadaili.com/"
}

func (s nimadaili) Run() {
	getProxy(s)
}

func (s nimadaili) Name() string {
	return "nimadaili"
}

func (s nimadaili) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
