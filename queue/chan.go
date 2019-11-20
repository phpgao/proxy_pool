package queue

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
)

var (
	config       = util.GetConfig()
	NewProxyChan = make(chan *model.HttpProxy, config.Concurrence)
	OldProxyChan = make(chan *model.HttpProxy, config.Concurrence)
)

func GetNewChan() chan *model.HttpProxy {
	return NewProxyChan
}

func GetOldChan() chan *model.HttpProxy {
	return OldProxyChan
}
