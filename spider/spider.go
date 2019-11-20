package spider

import (
	"github.com/phpgao/proxy_pool/model"
)

var spiderList []Crawler

// enabled spiders

func init() {
	spiderList = []Crawler{
		&IpHai{},
		&Rudnkh{},
		&CoolProxy{},
		&Xici{},
		&Spys{},
		&PubProxy{},
		&KuaiProxy{},
		&cn66{},
		&feiyi{},
		&ip89{},
		&goubanjia{},
		&jiangxianli{},
	}
}

func GetSpiders(ch chan<- *model.HttpProxy) []Crawler {
	for _, v := range spiderList {
		v.SetProxyChan(ch)
	}
	return spiderList
}
