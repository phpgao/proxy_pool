package source

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/parnurzeal/gorequest"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
	"time"
)

func (s *premProxy) StartUrl() []string {
	return []string{
		"https://premproxy.com/list/time-01.htm",
	}
}

func (s *premProxy) GetReferer() string {
	return "http://premproxy.com/"
}

type premProxy struct {
	Spider
}

func (s *premProxy) Cron() string {
	return "@every 7m"
}

func (s *premProxy) Name() string {
	return "premProxy"
}

func (s *premProxy) Run() {
	getProxy(s)
}

func (s *premProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	script := htmlquery.FindOne(doc, "/html/head/script[2]/@src")
	if script == nil {
		return
	}
	scriptUrl := htmlquery.InnerText(script)
	if scriptUrl == "" {
		return
	}
	scriptUrl = fmt.Sprintf("https://premproxy.com%s", scriptUrl)

	_, jsCode, errs := gorequest.New().Get(scriptUrl).
		Set("User-Agent", util.GetRandomUA()).
		Set("Content-Type", "text/html; charset=utf-8").
		Set("Referer", s.GetReferer()).
		Set("Pragma", `no-cache`).
		Timeout(time.Duration(s.TimeOut()) * time.Second).End()

	if len(errs) > 0 {
		err = errs[0]
		return
	}

	jsCode = strings.Replace(jsCode, "eval(", "a = (", 1)

	vm := otto.New()
	jsRs, err := vm.Run(jsCode)
	if err != nil {
		return
	}

	jsDecoded := jsRs.String()
	var cssPortMap = map[string]string{}

	allBlocks := regexp.MustCompile(`\$\('\.(\w+)'\)\.html\((\d+)\);`).FindAllStringSubmatch(jsDecoded, -1)

	for _, block := range allBlocks {
		cssPortMap[block[1]] = block[2]
	}

	proxyList := htmlquery.Find(doc, "//*[@id='proxylistt']/tbody/tr")
	for _, n := range proxyList {
		tmp := htmlquery.FindOne(n, "//td[1]/span/input/@value")
		value := htmlquery.InnerText(tmp)
		if !strings.Contains(value, "|") {
			continue
		}
		ipPort := strings.Split(value, "|")
		ip := strings.TrimSpace(ipPort[0])
		port := strings.TrimSpace(cssPortMap[ipPort[1]])

		proxies = append(proxies, &model.HttpProxy{
			Ip:   ip,
			Port: port,
		})
	}
	return
}
