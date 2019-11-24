package util

import (
	"github.com/ipipdotnet/ipdb-go"
)

var (
	Db *ipdb.City
)

func init() {
	var err error
	Db, err = ipdb.NewCity("ipipfree.ipdb")
	if err != nil {
		logger.WithError(err).Fatal("error initial ipdb")
	}
}
