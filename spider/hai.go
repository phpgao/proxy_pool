package spider

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/antchfx/htmlquery"
	"strings"
)

type IpHai struct {
	Spider
}

func (s *IpHai) Cron() string {
	return "@every 1m"
}

func (s *IpHai) StartUrl() []string {
	return []string{
		"http://www.iphai.com/free/ng",
		"http://www.iphai.com/free/np",
		"http://www.iphai.com/free/wg",
		"http://www.iphai.com/free/wp",
	}
}

func (s *IpHai) RandomDelay() bool {
	return true
}

func (s *IpHai) GetReferer() string {
	return "http://www.iphai.com/"
}

func (s *IpHai) Run() {
	getProxy(s)
}

func (s *IpHai) Name() string {
	return "ipHai"
}

func (s *IpHai) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	list := htmlquery.Find(doc, "/html/body/div[2]/div[2]/table/tbody/tr[position()>1]")

	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))
		schema := htmlquery.InnerText(htmlquery.FindOne(n, "//td[4]"))

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