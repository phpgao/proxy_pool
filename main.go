package main

import (
	"fmt"
	"github.com/phpgao/proxy_pool/schedule"
	"github.com/phpgao/proxy_pool/server"
	"github.com/phpgao/proxy_pool/util"
	"github.com/phpgao/proxy_pool/validator"
	"syscall"
)

var (
	logger = util.GetLogger()
)

const (
	VERSION = "0.2"
)

func main() {
	ShowWelcome()

	if util.ServerConf.Manager {
		logger.Info("Running in as Manager")

		scheduler := schedule.NewScheduler()
		go scheduler.Run()
	}
	if util.ServerConf.Worker {
		logger.Info("Running in as Worker")
		go validator.NewValidator()
		go validator.OldValidator()
	}

	go server.ServeReverse()
	server.Serve()

}

func ShowWelcome() {
	fmt.Printf(`
______                                        _ 
| ___ \                                      | |
| |_/ / __ _____  ___   _   _ __   ___   ___ | |        Proxy pool v%s
|  __/ '__/ _ \ \/ / | | | | '_ \ / _ \ / _ \| |        Proxy port: %d
| |  | | | (_) >  <| |_| | | |_) | (_) | (_) | |        Api port: %d
\_|  |_|  \___/_/\_\\__, | | .__/ \___/ \___/|_|        %s
                     __/ | | |                  
                    |___/  |_|                  `, VERSION, util.ServerConf.ApiPort, util.ServerConf.ProxyPort, "https://phpgao.com")
	fmt.Println()
	logger.Info("Proxy_pool is starting")
	logger.Infof("Proxy_pool VERSION == %s", VERSION)
	logger.Info("Configuration loaded")
	ulimit()
}

func ulimit() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logger.WithError(err).Error("Error Getting Rlimit")
	}
	rLimit.Max = uint64(util.ServerConf.UlimitMax)
	rLimit.Cur = uint64(util.ServerConf.UlimitCur)
	logger.WithField("ulimit", rLimit).Info("Try Setting Rlimit")

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logger.WithError(err).Error("Error Setting Rlimit")
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logger.WithError(err).Error("Error Getting Rlimit")
	}
	logger.WithField("ulimit", rLimit).Info("Rlimit Final")
}
