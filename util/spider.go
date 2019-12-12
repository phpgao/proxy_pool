package util

import (
	"encoding/json"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/parnurzeal/gorequest"
	"net"
	"regexp"
	"strings"
	"time"
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

func GetWsFromChrome(url string) (ws string, err error) {
	host, port := Parse(url)
	var chromeApi string
	// if ip format
	//if IsIpFormat(host) {
	//	chromeApi = fmt.Sprintf("http://%s:%s/json", host, port)
	//}
	//var chromeApi string
	//if IsIpFormat(host) {
	//	chromeApi = fmt.Sprintf("http://%s:%s/json", host, port)
	//} else {
	//	var addr []net.IP
	//	addr, err = net.LookupIP(host)
	//	if err != nil {
	//		return
	//	}
	//	host = addr[0].String()
	//	chromeApi = fmt.Sprintf("http://%s:%s/json", host, port)
	//}
	chromeApi = fmt.Sprintf("http://%s:%s/json", host, port)
	logger.WithField("chromeApi", chromeApi).Debug("get chromeApi")
	s := gorequest.New().Get(chromeApi).Timeout(5 * time.Second)

	if !IsIpFormat(host) {
		s.Header = map[string]string{
			"Host": "localhost",
		}
	}

	_, jsonBody, errs := s.End()
	if len(errs) > 0 {
		err = errs[0]
		return
	}
	var data []struct {
		Description          string `json:"description"`
		DevtoolsFrontendURL  string `json:"devtoolsFrontendUrl"`
		ID                   string `json:"id"`
		Title                string `json:"title"`
		Type                 string `json:"type"`
		URL                  string `json:"url"`
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}

	err = json.Unmarshal([]byte(jsonBody), &data)
	if err != nil {
		return
	}

	ws = strings.Replace(data[len(data)-1].WebSocketDebuggerURL, "localhost", fmt.Sprintf("%s:%s", host, port), 1)

	return
}

func Parse(url string) (string, string) {
	if strings.Contains(url, ":") {
		t := strings.Split(url, ":")
		return t[0], t[1]
	}

	return url, "9222"
}

func IsIpFormat(ip string) bool {
	address := net.ParseIP(ip)
	if address == nil {
		return false
	}

	return true
}
