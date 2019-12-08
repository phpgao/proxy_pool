package job

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"github.com/robertkrimen/otto"
	"regexp"
	"strconv"
	"strings"
)

type ultraProxies struct {
	Spider
}

func (s *ultraProxies) StartUrl() []string {
	return []string{
		"http://www.ultraproxies.com/",
	}
}

func (s *ultraProxies) Cron() string {
	return "@every 1h"
}

func (s *ultraProxies) Name() string {
	return "ultraProxies"
}

func (s *ultraProxies) GetReferer() string {
	return "http://www.ultraproxies.com/"
}

func (s *ultraProxies) Run() {
	getProxy(s)
}

func (s *ultraProxies) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	scriptNode := htmlquery.FindOne(doc, "//*[@id='inner']/script")

	if scriptNode == nil {
		err = errors.New("no script found")
		return
	}

	jsCode := htmlquery.InnerText(scriptNode)

	jsCode = regexp.MustCompile(`\$\("\.port.*`).ReplaceAllString(jsCode, "")

	vm := otto.New()
	xValue, err := vm.Run(jsCode)
	if err != nil {
		return
	}
	Value, err := xValue.ToInteger()
	if err != nil {
		return
	}
	proxyList := htmlquery.Find(doc, "//*[@id='inner']/table/tbody/tr/td[3]/table[2]/tbody/tr[position()>1]")

	for _, n := range proxyList {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))

		ip = strings.TrimSpace(ip)
		ip = strings.TrimRight(ip, ":")
		port, err = decodePort(int(Value), strings.TrimSpace(port))
		if err != nil {
			continue
		}
		proxies = append(proxies, &model.HttpProxy{
			Ip:   ip,
			Port: port,
		})
	}

	return
}

func decodePort(x int, port string) (portStr string, err error) {
	asciiPorts := strings.Split(port, "-")
	for _, one := range asciiPorts {
		var tmpInt int
		tmpInt, err = strconv.Atoi(one)
		if err != nil {
			return
		}

		portStr += string(rune(x + tmpInt))
	}

	return
}
