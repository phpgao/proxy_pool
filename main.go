package main

import (
	"github.com/phpgao/proxy_pool/schedule"
	"github.com/phpgao/proxy_pool/server"
	"github.com/phpgao/proxy_pool/util"
	"github.com/phpgao/proxy_pool/validator"
)

var (
	config = util.GetConfig()
	//logger = util.GetLogger()
)

func main() {
	if config.Manager {
		scheduler := schedule.NewScheduler()
		go scheduler.Run()
	}
	if config.Worker {
		go validator.NewValidator()
		go validator.OldValidator()
	}
	go server.ServeReverse()
	server.Serve()
}
