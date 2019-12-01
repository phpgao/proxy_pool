// +build !windows

package ulimit

import (
	"github.com/phpgao/proxy_pool/util"
	"syscall"
)

var (
	logger = util.GetLogger()
)

func Set() {
	var rLimit syscall.Rlimit

	rLimit.Max = uint64(util.ServerConf.UlimitMax)
	rLimit.Cur = uint64(util.ServerConf.UlimitCur)
	logger.WithField("ulimit", rLimit).Info("Try Setting Rlimit")

	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logger.WithError(err).Info("Error Setting Rlimit")
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logger.WithError(err).Info("Error Getting Rlimit")
	}
	logger.WithField("ulimit", rLimit).Info("Rlimit Final")
}
