package job

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"net/url"
	"strings"
)

type proxyLists struct {
	Spider
}

func (s *proxyLists) StartUrl() []string {
	return []string{
		"http://www.proxylists.net/cn_0_ext.html",
		"http://www.proxylists.net/cn_1_ext.html",
		"http://www.proxylists.net/cn_2_ext.html",
		"http://www.proxylists.net/cn_3_ext.html",
		"http://www.proxylists.net/cn_4_ext.html",
		"http://www.proxylists.net/cn_5_ext.html",
		"http://www.proxylists.net/us_0_ext.html",
		"http://www.proxylists.net/gb_0_ext.html",
		"http://www.proxylists.net/ca_0_ext.html",
		"http://www.proxylists.net/au_0_ext.html",
	}
}

func (s *proxyLists) Cron() string {
	return "@every 5m"
}

func (s *proxyLists) Name() string {
	return "proxy_lists"
}

func (s *proxyLists) GetReferer() string {
	return "http://www.proxyLists.com/"
}

func (s *proxyLists) Run() {
	getProxy(s)
}

func (s *proxyLists) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	proxyList := htmlquery.Find(doc, "/html/body/font/b/table/tbody/tr[1]/td[2]/table/tbody/tr[position()>2 and position()<last()]")
	for _, n := range proxyList {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		if strings.Contains(ip, "Page") {
			continue
		}
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))
		proxyType := htmlquery.InnerText(htmlquery.FindOne(n, "//td[3]"))
		if proxyType == "Anonymous" || proxyType == "Transparent" {
			ip = strings.TrimSpace(ip)
			ip = s.decodeIp(strings.TrimSpace(ip))
			if port != "" {
				proxies = append(proxies, &model.HttpProxy{
					Ip:   ip,
					Port: port,
				})
			}
		}
	}

	return
}

func (s proxyLists) decodeIp(port string) string {
	tmp, err := url.QueryUnescape(port)
	if err != nil {
		return ""
	}

	return util.FindIp(tmp)
}
