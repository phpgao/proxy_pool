package spider

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/antchfx/htmlquery"
	"strings"
)

type Xici struct {
	Spider
}

func (s *Xici) Cron() string {
	return "@every 1m"
}

func (s *Xici) GetReferer() string {
	return "https://www.xicidaili.com/"
}

func (s *Xici) StartUrl() []string {
	return []string{
		"http://www.xicidaili.com/nn",
		"http://www.xicidaili.com/wn",
	}
}

func (s *Xici) Run() {
	getProxy(s)
}

func (s *Xici) Name() string {
	return "Xici"
}

func (s *Xici) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	list := htmlquery.Find(doc, "//*[@id='ip_list']/tbody/tr[position()>1]")
	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[3]"))
		schema := htmlquery.InnerText(htmlquery.FindOne(n, "//td[6]"))

		ip = strings.TrimSpace(ip)
		port = strings.TrimSpace(port)
		schema = strings.TrimSpace(schema)
		if len(schema) == 0 {
			schema = "http"
		}

		proxies = append(proxies, &model.HttpProxy{
			Ip:        ip,
			Port:      port,
			Schema:    strings.ToLower(schema),
			Score:     config.Score,
			Latency:   0,
			From:      s.Name(),
			Anonymous: 0,
		})

	}
	return
}
