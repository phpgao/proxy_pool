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
		"https://www.zdaye.com/FreeIPList.html",
	}
}

func (s *zdy) Enabled() bool {
	return true
}
func (s *zdy) Cron() string {
	return "@every 5m"
}

func (s *zdy) GetReferer() string {
	return "https://www.zdaye.com"
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

	list := htmlquery.Find(doc, "//table/tbody/tr[position()>1]")
	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		ip = strings.TrimSpace(ip)
		for _, p := range []string{
			"80", "8080", "8008", "8888", "9999", "1080", "3000",
		} {
			proxies = append(proxies, &model.HttpProxy{
				Ip:   ip,
				Port: p,
			})
		}

	}
	return
}
