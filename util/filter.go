package util

import (
	"github.com/phpgao/proxy_pool/model"
	"net"
	"strconv"
	"strings"
)

func filterOfSchema(v string) func(model.HttpProxy) bool {
	return func(proxy model.HttpProxy) bool {
		if proxy.Schema == strings.ToLower(v) {
			return true
		}
		return false
	}
}

func filterOfCn(v string) func(model.HttpProxy) bool {
	return func(proxy model.HttpProxy) bool {
		if proxy.Country == strings.ToLower(v) {
			return true
		}
		return false
	}
}
func filterOfScore(v int) func(model.HttpProxy) bool {
	return func(proxy model.HttpProxy) bool {
		if proxy.Score >= v {
			return true
		}
		return false
	}
}
func filterOfSource(v string) func(model.HttpProxy) bool {
	return func(proxy model.HttpProxy) bool {
		if strings.ToLower(proxy.From) == strings.ToLower(v) {
			return true
		}
		return false
	}
}

func GetNewFilter(options map[string]string) (f []func(model.HttpProxy) bool, err error) {
	for k, v := range options {
		if k == "schema" && v != "" {
			f = append(f, filterOfSchema(v))
		}
		if k == "source" && v != "" {
			f = append(f, filterOfSource(v))
		}
		if k == "score" && v != "" {
			i := 0
			i, numError := strconv.Atoi(v)
			if numError != nil {
				if _, ok := numError.(*strconv.NumError); ok {
					i = 0
				} else {
					err = numError
					return
				}
			}
			f = append(f, filterOfScore(i))
		}

		if k == "country" && v != "" {
			f = append(f, filterOfCn(v))
		}
	}

	return
}

func FilterProxy(proxy *model.HttpProxy) bool {
	if tmp := net.ParseIP(proxy.Ip); tmp.To4() == nil {
		return false
	}

	port, err := strconv.Atoi(proxy.Port)
	if err != nil {
		return false
	}

	if port < 1 || port > 65535 {
		return false
	}

	ipInfo, err := Db.FindInfo(proxy.Ip, "CN")
	if err != nil {
		logger.WithField("ip", proxy.Ip).WithError(err).Warn("can not find ip info")
		return false
	}

	if ipInfo.CountryName == "中国" {
		proxy.Country = "cn"
	} else {
		proxy.Country = ipInfo.CountryName
	}

	return true
}
