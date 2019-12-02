package ipdb

import (
	"github.com/phpgao/proxy_pool/util"
)

var (
	Db     *City
	logger = util.GetLogger()
)

func init() {
	var err error
	r, e := newReaderFromGo("ipipfree.ipdb", &CityInfo{})
	if e != nil {
		panic(e)
	}

	Db = &City{
		reader: r,
	}
	if err != nil {
		logger.WithError(err).Fatal("error initial ipdb")
	}
}
