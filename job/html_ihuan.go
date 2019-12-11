package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type ihuan struct {
	Spider
}

func (s *ihuan) StartUrl() []string {
	return []string{
		"https://ip.ihuan.me/ti.html",
		"https://ip.ihuan.me/",
		"https://ip.ihuan.me/address/5Lit5Zu9.html",
	}
}

func (s *ihuan) Cron() string {
	return "@every 30m"
}

func (s *ihuan) GetReferer() string {
	return "https://ip.ihuan.me/"
}

func (s *ihuan) Run() {
	getProxy(s)
}

func (s *ihuan) Name() string {
	return "ihuan"
}

func (s *ihuan) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
