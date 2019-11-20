package spider

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/antchfx/htmlquery"
	"strings"
)

func (s *feiyi) StartUrl() []string {
	return []string{
		"http://www.feiyiproxy.com/?page_id=1457",
	}
}

func (s *feiyi) GetReferer() string {
	return "http://www.feiyiproxy.com/"
}

type feiyi struct {
	Spider
}

func (s *feiyi) Cron() string {
	return "@every 5m"
}

func (s *feiyi) Name() string {
	return "feiyi"
}

func (s *feiyi) Run() {
	getProxy(s)
}

func (s *feiyi) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	list := htmlquery.Find(doc, "//*[@id='post-1457']/div/div/div[3]/div/div/div/table/tbody/tr[position()>1]")
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
		if anonymous == "高匿" || anonymous == "普匿"|| anonymous == "透明"{
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
