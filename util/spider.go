package util

import (
	"github.com/corpix/uarand"
	"regexp"
)

const (
	RegIp                = `(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])`
	RegProxy             = `(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5]):\d{0,5}`
	RegProxyWithoutColon = `(?:(?:[0,1]?\d?\d|2[0-4]\d|25[0-5])\.){3}(?:[0,1]?\d?\d|2[0-4]\d|25[0-5]) \d{0,5}`
)

func GetRandomUA() string {
	return uarand.GetRandom()
}

func FindIp(s string) string {
	reg := regexp.MustCompile(RegIp)
	rs := reg.FindAllString(s, 1)
	if rs == nil {
		return ""
	}

	return rs[0]
}
