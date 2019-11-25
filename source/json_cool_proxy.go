package source

import (
	"encoding/json"
	"github.com/apex/log"
	"sort"
	"strconv"

	"github.com/phpgao/proxy_pool/model"
)

type coolProxy struct {
	Spider
}

func (s *coolProxy) StartUrl() []string {
	return []string{
		"https://cool-proxy.net/proxies.json",
	}
}

func (s *coolProxy) GetReferer() string {
	return "https://cool-proxy.net"
}

func (s *coolProxy) Run() {
	getProxy(s)
}

func (s *coolProxy) Cron() string {
	return "@every 5m"
}

func (s *coolProxy) Name() string {
	return "cool_proxy"
}

func (s *coolProxy) TimeOut() int {
	return 30
}

type coolProxyJson struct {
	ResponseTimeAverage  float64 `json:"response_time_average"`
	IP                   string  `json:"ip"`
	Score                float64 `json:"score"`
	UpdateTime           float64 `json:"update_time"`
	WorkingAverage       float64 `json:"working_average"`
	Anonymous            int     `json:"anonymous"`
	CountryCode          string  `json:"country_code"`
	CountryName          string  `json:"country_name"`
	Port                 int     `json:"port"`
	DownloadSpeedAverage float64 `json:"download_speed_average"`
}

func (s *coolProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	var coolProxies []coolProxyJson
	err = json.Unmarshal([]byte(body), &coolProxies)
	if err != nil {
		logger.WithError(err).WithFields(log.Fields{
			"body":    body,
			"timeout": s.TimeOut(),
		}).Debug("error parse json")
		return
	}

	if len(coolProxies) == 0 {
		return
	}
	//sort by time desc
	sort.Slice(coolProxies, func(i, j int) bool {
		return coolProxies[i].UpdateTime > coolProxies[j].UpdateTime
	})

	for _, proxy := range coolProxies {
		proxies = append(proxies, &model.HttpProxy{
			Ip:        proxy.IP,
			Port:      strconv.Itoa(proxy.Port),
			Anonymous: proxy.Anonymous,
		})
	}
	return
}
