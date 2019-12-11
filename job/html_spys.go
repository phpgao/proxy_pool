package job

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/apex/log"
	"github.com/parnurzeal/gorequest"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"github.com/robertkrimen/otto"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type spys struct {
	Spider
}

func (s *spys) StartUrl() []string {
	return []string{
		"http://spys.one/en/anonymous-proxy-list/",
		"http://spys.one/free-proxy-list/CHN/",
		"http://spys.one/free-proxy-list/US/",
	}
}

func (s *spys) Cron() string {
	return "@every 5m"
}

func (s *spys) Name() string {
	return "spys"
}

func (s *spys) GetReferer() string {
	return "http://spys.one/en/anonymous-proxy-list/"
}

func (s *spys) Run() {
	getProxy(s)
}

func (s *spys) Fetch(siteUrl string, useProxy bool) (body string, err error) {

	if s.RandomDelay() {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}

	request := gorequest.New()
	var superAgent *gorequest.SuperAgent
	var resp gorequest.Response
	var errs []error
	superAgent = request.Post(siteUrl).
		Set("User-Agent", util.GetRandomUA()).
		Set("Content-Type", `text/html; charset=utf-8`).
		Set("Referer", s.GetReferer()).
		Set("Pragma", `no-cache`).
		Send("xpp=2&xf1=1&xf2=0&xf4=0&xf5=1").
		Timeout(time.Duration(s.TimeOut()) * time.Second).SetDebug(util.ServerConf.DumpHttp)

	if useProxy {
		var proxy model.HttpProxy
		proxy, err = storeEngine.Random()
		if err != nil {
			return
		}
		p := fmt.Sprintf("http://%s:%s", proxy.Ip, proxy.Port)
		logger.WithFields(log.Fields{
			"proxy": p,
			"url":   siteUrl,
		}).Debug("with proxy")
		resp, body, errs = superAgent.Proxy(p).End()
	} else {
		resp, body, errs = superAgent.End()
	}
	if err = s.errAndStatus(errs, resp); err != nil {
		return
	}

	body = strings.TrimSpace(body)
	return
}

func (s *spys) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}
	list := htmlquery.Find(doc, "/html/body/table[2]/tbody/tr[5]/td/table/tbody/tr[@onmouseover]")
	var initJs string
	initJsBlock := htmlquery.Find(doc, "/html/body/script")
	for _, script := range initJsBlock {
		initJs = htmlquery.InnerText(script)
	}
	var vm *otto.Otto
	if initJs != "" {
		vm = otto.New()
		_, err = vm.Run(initJs)
		if err != nil {
			return
		}
	}

	for _, n := range list {
		ipText := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
		matchedIp := re.FindAllString(ipText, -1)
		if len(matchedIp) > 0 {
			portJs := getPortJs(ipText)
			port, err := ParsePort(vm, portJs)
			if err != nil {
				continue
			}
			proxies = append(proxies, &model.HttpProxy{
				Ip:   matchedIp[0],
				Port: port,
			})
		}
	}
	return
}

func getPortJs(s string) (js string) {
	sL := len(s)
	match := `<font class=spy2>:<\/font>"`
	i := strings.Index(s, `<font class=spy2>:<\/font>"+`)
	l := len(match) + i
	return s[l : sL-1]
}

func ParsePort(vm *otto.Otto, PortJs string) (port string, err error) {
	code := fmt.Sprintf("\"\"%s", PortJs)
	value, err := vm.Run(code)
	if err != nil {
		return
	}
	port, err = value.ToString()
	if err != nil {
		return
	}
	return
}
