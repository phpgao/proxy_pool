package util

import (
	"fmt"
	"github.com/jinzhu/configor"
)

var c Config

type Config struct {
	Worker  bool `default:"true"`
	Manager bool `default:"true"`
	Redis   struct {
		IP        string `default:"127.0.0.1"`
		Port      int    `default:"6379"`
		Db        int    `default:"1"`
		Auth      string `default:""`
		PrefixKey string `default:"proxy_pool"`
	}
	Concurrence   int  `default:"100"`
	Debug         bool `default:"false"`
	Timeout       int  `default:"10"`
	CheckInterval int  `default:"60"`
	Port          int  `default:"8080"`
	Init          bool `default:"false"`
	Expire        int  `default:"0"`
	Score         int  `default:"60"`
}

func (c Config) GetInternalCron() string {
	return fmt.Sprintf("@every %ds", c.CheckInterval)
}

func GetConfig() *Config {
	if !c.Init {
		err := configor.Load(&c, "config.json")
		if err != nil {
			panic(err)
		}
		c.Init = true
	}
	return &c
}
