package source

import (
	"encoding/base64"
	"github.com/phpgao/proxy_pool/model"
	"regexp"
	"strings"
)

func (s *proxylist) StartUrl() []string {
	return []string{
		"http://proxy-list.org/english/index.php",
	}
}

func (s *proxylist) GetReferer() string {
	return "http://proxy-list.org"
}

type proxylist struct {
	Spider
}

func (s *proxylist) Cron() string {
	return "@every 2m"
}

func (s *proxylist) Name() string {
	return "proxylist"
}

func (s *proxylist) Run() {
	getProxy(s)
}

func (s *proxylist) Parse(body string) (proxies []*model.HttpProxy, err error) {
	reg := regexp.MustCompile(`Proxy\('([a-zA-Z0-9/+]*={0,2})'\)`)
	rs := reg.FindAllStringSubmatch(body, -1)
	for _, proxy := range rs {
		proxyByte, err := base64.StdEncoding.DecodeString(proxy[1])
		if err != nil {
			logger.WithError(err).Debug("error decode")
			continue
		}
		proxyStr := string(proxyByte)
		if strings.Contains(proxyStr, ":") {
			proxyInfo := strings.Split(proxyStr, ":")
			proxies = append(proxies, &model.HttpProxy{
				Ip:   proxyInfo[0],
				Port: proxyInfo[1],
			})
		}
	}
	return
}
