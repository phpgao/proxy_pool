package job

import (
	"github.com/phpgao/proxy_pool/model"
)

var ListOfSpider = []Crawler{
	&IpHai{},
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
	&freeip{},
	&ab57{},
	&clarketm{},
	&httptunnel{},
	&proxylist{},
	&proxylistplus{},
	&aliveProxy{},
	&proxyDb{},
	&usProxy{},
	&siteDigger{},
	&dogdev{},
	&newProxy{},
	&xseo{},
	&ultraProxies{},
	&premProxy{},
	&nntime{},
	&proxyListsLine{},
	&myProxy{},
	&proxyIpList{},
	&blackHat{},
	&proxyLists{},
	&ip3366{},
	&xiladaili{},
	&nimadaili{},
	&zdy{},
}

func GetSpiders(ch chan<- *model.HttpProxy) []Crawler {
	for _, v := range ListOfSpider {
		v.SetProxyChan(ch)
	}
	return ListOfSpider
}
