package job

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/phpgao/proxy_pool/model"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

type proxyDb struct {
	Spider
}

func (s *proxyDb) StartUrl() []string {
	return []string{
		"http://proxydb.net/?protocol=http&protocol=https&country=",
	}
}

func (s *proxyDb) Cron() string {
	return "@every 1h"
}

func (s *proxyDb) Name() string {
	return "proxydb"
}

func (s *proxyDb) GetReferer() string {
	return "http://proxydb.net/"
}

func (s *proxyDb) Run() {
	getProxy(s)
}

func (s *proxyDb) Parse(body string) (proxies []*model.HttpProxy, err error) {
	scriptRe := regexp.MustCompile(`<div style="display:none" data-(\w+)+="(\d+)"></div>`)
	scriptRs := scriptRe.FindAllStringSubmatch(body, 1)
	if scriptRs == nil {
		err = errors.New("random data not found")
		return
	}

	randomKey := scriptRs[0][1]
	randomValue := scriptRs[0][2]

	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	list := htmlquery.Find(doc, "//td/script")
	vm := otto.New()
	err = registerJsFunctions(vm)
	if err != nil {
		return
	}
	for _, n := range list {
		scriptString := htmlquery.InnerText(n)
		old := fmt.Sprintf("(+document.querySelector('[data-%s]').getAttribute('data-%s'))", randomKey, randomKey)

		scriptString = strings.Replace(scriptString, old, fmt.Sprintf("\"%s\"", randomValue), 1)
		scriptString = regexp.MustCompile(`document.write\('<a href="/' \+ (\w)\s.*;`).ReplaceAllString(scriptString, "proxy = $1 + yxy + String.fromCharCode(58) + pp")

		proxy, err := DecodeProxy(vm, scriptString)
		if err != nil {
			continue
		}
		if strings.Contains(proxy, ":") {
			proxyInfo := strings.Split(proxy, ":")
			proxies = append(proxies, &model.HttpProxy{
				Ip:   proxyInfo[0],
				Port: proxyInfo[1],
			})
		}
	}

	return
}

func DecodeProxy(vm *otto.Otto, js string) (proxy string, err error) {
	_, err = vm.Run(js)
	if err != nil {
		return
	}
	value, err := vm.Get("proxy")
	if err != nil {
		return
	}
	proxy, err = value.ToString()
	if err != nil {
		return
	}
	return
}

func jsBtoa(b string) string {
	return base64.StdEncoding.EncodeToString([]byte(b))
}

func jsAtob(str string) string {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		v, _ := otto.ToValue("Failed to execute 'jsAtob': The string to be decoded is not correctly encoded.")
		panic(v)
	}
	return string(b)
}

func registerJsFunctions(vm *otto.Otto) (err error) {
	err = vm.Set("btoa", jsBtoa)
	if err != nil {
		return err
	}
	err = vm.Set("atob", jsAtob)
	if err != nil {
		return err
	}
	return
}
