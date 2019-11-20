package spider

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/phpgao/proxy_pool/model"
)

type CoolProxy struct {
	Spider
}

func (s *CoolProxy) StartUrl() []string {
	return []string{
		"https://cool-proxy.net/proxies.json",
	}
}

func (s *CoolProxy) GetReferer() string {
	return "https://cool-proxy.net"
}

func (s *CoolProxy) Run() {
	getProxy(s)
}

func (s *CoolProxy) Cron() string {
	return "@every 1m"
}

func (s *CoolProxy) Name() string {
	return "cool_proxy"
}

type coolProxy struct {
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

func (s *CoolProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	var coolProxies []coolProxy
	err = json.Unmarshal([]byte(body), &coolProxies)
	if err != nil {
		logger.WithError(err).WithField("body", body).Debug("error parse json")
		return
	}

	if len(coolProxies) == 0 {
		return
	}
	//sort by time desc
	sort.Slice(coolProxies, func(i, j int) bool {
		return coolProxies[i].UpdateTime > coolProxies[j].UpdateTime
	})

	for _, proxy := range coolProxies[:100] {
		proxies = append(proxies, &model.HttpProxy{
			Ip:        proxy.IP,
			Port:      strconv.Itoa(proxy.Port),
			Schema:    "http",
			Score:     config.Score,
			Latency:   0,
			From:      s.Name(),
			Anonymous: proxy.Anonymous,
		})
	}
	return
}
