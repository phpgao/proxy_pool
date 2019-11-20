package schedule

import (
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
)

type Internal struct {
	channel chan *model.HttpProxy
	db      db.Store
}

func (i Internal) Run() {
	m := i.db.GetAll()
	proxyCount := len(m)
	logger.WithField("Count", proxyCount).Info("start internal check")
	for k := range m {
		i.channel <- &m[k]
	}
}
