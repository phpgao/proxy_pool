package queue

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
)

var (
	config       = util.ServerConf
	NewProxyChan = make(chan *model.HttpProxy, config.NewQueue)
	OldProxyChan = make(chan *model.HttpProxy, config.OldQueue)
)

func GetNewChan() chan *model.HttpProxy {
	return NewProxyChan
}

func GetOldChan() chan *model.HttpProxy {
	return OldProxyChan
}
