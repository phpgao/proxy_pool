package util

import (
	"encoding/json"
	"github.com/corpix/uarand"
	"github.com/parnurzeal/gorequest"
	"regexp"
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
	_, jsonBody, errs := gorequest.New().Get(url).Timeout(5 * time.Second).End()
	if len(errs) > 0 {
		err = errs[0]
		panic(err)
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
		panic(err)
	}
	ws = data[len(data)-1].WebSocketDebuggerURL
	return
}
