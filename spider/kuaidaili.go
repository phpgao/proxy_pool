package spider

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type KuaiProxy struct {
	Spider
}

func (s *KuaiProxy) StartUrl() []string {
	return []string{
		"https://www.kuaidaili.com/free/intr/",
		"https://www.kuaidaili.com/free/inha/",
	}
}

func (s *KuaiProxy) Cron() string {
	return "@every 1m"
}

func (s *KuaiProxy) GetReferer() string {
	return "https://www.kuaidaili.com/"
}

func (s *KuaiProxy) Run() {
	getProxy(s)
}

func (s *KuaiProxy) Name() string {
	return "Kuai"
}

func (s *KuaiProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	list := htmlquery.Find(doc, "//*[@id='list']/table/tbody/tr[position()>1]")
	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))
		schema := htmlquery.InnerText(htmlquery.FindOne(n, "//td[4]"))
		anonymous := htmlquery.InnerText(htmlquery.FindOne(n, "//td[3]"))

		ip = strings.TrimSpace(ip)
		port = strings.TrimSpace(port)
		schema = strings.TrimSpace(schema)
		anonymous = strings.TrimSpace(anonymous)
		var anonymousInt int
		if anonymous == "高匿名" {
			anonymousInt = 1
		} else {
			anonymousInt = 0
		}
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
			Anonymous: anonymousInt,
		})
	}
	return
}
