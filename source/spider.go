package source

import (
	"github.com/phpgao/proxy_pool/model"
)

var ListOfSpider = []Crawler{
	//&IpHai{},
	&rudnkh{},
	&coolProxy{},
	&xici{},
	&spys{},
	&pubProxy{},
	&kuaiProxy{},
	&cn66{},
	&feiyi{},
	&ip89{},
	&goubanjia{},
	&jiangxianli{},
	&ab57{},
	&clarketm{},
	&httptunnel{},
}

func GetSpiders(ch chan<- *model.HttpProxy) []Crawler {
	for _, v := range ListOfSpider {
		v.SetProxyChan(ch)
	}
	return ListOfSpider
}