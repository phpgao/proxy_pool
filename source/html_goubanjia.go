package source

import (
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"strconv"
	"strings"
)

func (s *goubanjia) StartUrl() []string {
	return []string{
		"http://www.goubanjia.com/",
	}
}

func (s *goubanjia) GetReferer() string {
	return "http://www.goubanjiaproxy.com/"
}

type goubanjia struct {
	Spider
}

func (s *goubanjia) Cron() string {
	return "@every 1m"
}

func (s *goubanjia) Name() string {
	return "goubanjia"
}

func (s *goubanjia) Run() {
	getProxy(s)
}

func (s *goubanjia) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	list := htmlquery.Find(doc, "//td[@class='ip']")
	for _, n := range list {
		ipCss := htmlquery.Find(n, `.//*[not(
contains(@style, 'display:none;')
) and not(
contains(@style, 'display: none;')
) and not(
contains(@class, 'port')
)]/text()`)
		var tempString []string
		for _, i := range ipCss {
			x := htmlquery.InnerText(i)
			if x == ":" {
				continue
			}
			tempString = append(tempString, x)
		}
		ip := strings.Join(tempString, "")
		portCss := htmlquery.FindOne(n, `.//span[contains(@class, 'port')]/@class`)

		portString := htmlquery.InnerText(portCss)
		portString = strings.Replace(portString, "port ", "", 1)
		l := len(portString)
		portInt := 0
		str := portString
		for i := 0; i < l; i++ {
			portInt *= 10
			portInt += int(str[i]) - int('A')
		}
		portInt /= 8
		port := strconv.Itoa(portInt)

		proxies = append(proxies, &model.HttpProxy{
			Ip:   ip,
			Port: port,
		})
	}
	return
}
