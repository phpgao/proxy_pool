package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strings"
)

type xiladaili struct {
	Spider
}

func (s *xiladaili) StartUrl() []string {
	return []string{
		"http://www.xiladaili.com/gaoni/",
		"http://www.xiladaili.com/gaoni/2/",
		"http://www.xiladaili.com/gaoni/3/",
		"http://www.xiladaili.com/http/",
		"http://www.xiladaili.com/http/2/",
		"http://www.xiladaili.com/http/3/",
		"http://www.xiladaili.com/https/",
		"http://www.xiladaili.com/https/2/",
		"http://www.xiladaili.com/https/3/",
		"http://www.xiladaili.com/putong/",
		"http://www.xiladaili.com/putong/2/",
		"http://www.xiladaili.com/putong/3/",
	}
}

func (s *xiladaili) Cron() string {
	return "@every 2m"
}

func (s *xiladaili) GetReferer() string {
	return "http://www.xiladaili.com/"
}

func (s *xiladaili) Run() {
	getProxy(s)
}

func (s *xiladaili) Name() string {
	return "xiladaili"
}

func (s *xiladaili) Parse(body string) (proxies []*model.HttpProxy, err error) {
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
