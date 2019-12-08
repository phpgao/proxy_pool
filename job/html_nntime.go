package job

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

func (s *nntime) StartUrl() []string {
	return []string{
		"http://nntime.com/proxy-updated-01.htm",
	}
}

func (s *nntime) GetReferer() string {
	return "http://nntime.com/"
}

type nntime struct {
	Spider
}

func (s *nntime) Cron() string {
	return "@every 30m"
}

func (s *nntime) Name() string {
	return "nntime"
}

func (s *nntime) Run() {
	getProxy(s)
}

func (s *nntime) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	script := htmlquery.FindOne(doc, "/html/head/script[2]")
	if script == nil {
		return
	}
	jsCode := htmlquery.InnerText(script)
	if jsCode == "" {
		return
	}
	vm := otto.New()
	_, err = vm.Run(jsCode)
	if err != nil {
		return
	}

	proxyList := htmlquery.Find(doc, `//*[@id="proxylist"]/tbody/tr`)
	for _, n := range proxyList {
		tmp := htmlquery.FindOne(n, "//td[2]")
		value := htmlquery.InnerText(tmp)
		ip, err := getIp(value)
		if err != nil {
			continue
		}
		port, err := getPort(value, vm)
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

func getIp(s string) (ip string, err error) {
	r := regexp.MustCompile(util.RegIp).FindAllStringSubmatch(s, 1)
	if r == nil {
		err = errors.New("no ip found")
		return

	}
	ip = r[0][0]
	return
}

func getPort(s string, vm *otto.Otto) (ip string, err error) {
	match := regexp.MustCompile(`(document.write\(.*\))`).FindAllStringSubmatch(s, 1)
	if match == nil {
		err = errors.New("no port js found")
		return
	}
	jsCode := match[0][0]

	rs, err := vm.Run(strings.Replace(jsCode, "document.write(\":\"", "a=(\"\"", 1))
	if err != nil {
		return
	}

	ip, err = rs.ToString()
	if err != nil {
		return
	}
	return
}
