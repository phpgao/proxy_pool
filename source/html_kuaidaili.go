package source

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type kuaiProxy struct {
	Spider
}

func (s *kuaiProxy) StartUrl() []string {
	return []string{
		"https://www.kuaidaili.com/free/intr/",
		"https://www.kuaidaili.com/free/inha/",
	}
}

func (s *kuaiProxy) Cron() string {
	return "@every 30m"
}

func (s *kuaiProxy) GetReferer() string {
	return "https://www.kuaidaili.com/"
}

func (s *kuaiProxy) Run() {
	getProxy(s)
}

func (s *kuaiProxy) Name() string {
	return "Kuai"
}

func (s *kuaiProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	list := htmlquery.Find(doc, "//*[@id='list']/table/tbody/tr[position()>1]")
	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))
		anonymous := htmlquery.InnerText(htmlquery.FindOne(n, "//td[3]"))

		ip = strings.TrimSpace(ip)
		port = strings.TrimSpace(port)

		anonymous = strings.TrimSpace(anonymous)
		var anonymousInt int
		if anonymous == "高匿名" {
			anonymousInt = 1
		} else {
			anonymousInt = 0
		}

		proxies = append(proxies, &model.HttpProxy{
			Ip:        ip,
			Port:      port,
			Anonymous: anonymousInt,
		})
	}
	return
}
