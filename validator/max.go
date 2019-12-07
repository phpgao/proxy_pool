package validator

import (
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"sort"
	"sync"
)

var logger = util.GetLogger("max")

type Proxies []model.HttpProxy

func (p Proxies) Len() int {
	return len(p)
}

func (p Proxies) Less(i, j int) bool {
	return p[i].Score < p[j].Score
}

func (p Proxies) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

var (
	CurrentProxyCount int
	lock              sync.RWMutex
)

func Update() {
	logger.Debug("begin wash db")
	lock.Lock()
	defer lock.Unlock()
	CurrentProxyCount = db.GetDb().Len()
	if CurrentProxyCount > util.ServerConf.MaxProxy {
		logger.Debug("need wash db")
		Wash()
	}
	logger.Debug("done")
}

func CanDo() bool {
	lock.RLock()
	defer lock.RUnlock()
	return CurrentProxyCount < util.ServerConf.MaxProxy
}

// 删除低分代理，直到剩余指定个数
func Wash() {
	max := util.ServerConf.MaxProxy
	if max == 0 {
		return
	}
	all := db.GetDb().GetAll()
	l := len(all)
	if l <= max {
		return
	}

	tmp := Proxies(all)
	sort.Sort(tmp)
	err := db.GetDb().RemoveAll(tmp[:l-max])

	if err != nil {
		logger.WithError(err).Error("error wash db")
		return
	}
	logger.WithField("count", l-max).Debug("db washed")
	return
}
