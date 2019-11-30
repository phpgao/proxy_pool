package util

import (
	"fmt"
	"github.com/koding/multiconfig"
	"time"
)

var ServerConf *Config

type Config struct {
	Worker              bool   `default:"true"`
	Manager             bool   `default:"true"`
	IP                  string `default:"127.0.0.1"`
	Port                int    `default:"6379"`
	Db                  int    `default:"1"`
	Auth                string `default:""`
	PrefixKey           string `default:"proxy_pool"`
	NewQueue            int    `default:"200"`
	OldQueue            int    `default:"300"`
	Debug               bool   `default:"false"`
	Timeout             int    `default:"10"`
	CheckInterval       int    `default:"60"`
	Init                bool   `default:"false"`
	Expire              int    `default:"0"`
	Score               int    `default:"60"`
	Retry               int    `default:"3"`
	TcpTimeout          int    `default:"5"`
	ProxyTimeout        int    `default:"5"`
	ApiBind             string `default:"0.0.0.0"`
	ApiPort             int    `default:"8088"`
	ProxyBind           string `default:"0.0.0.0"`
	ProxyPort           int    `default:"8089"`
	TcpTestTimeOut      int    `default:"5"`
	HttpsConnectTimeOut int    `default:"5"`
	OnlyChina           bool   `default:"true"`
}

func init() {
	m := multiconfig.NewWithPath("config.json")
	serverConf := new(Config)
	m.MustLoad(serverConf)
	ServerConf = serverConf
}

func (c Config) GetInternalCron() string {
	return fmt.Sprintf("@every %ds", c.CheckInterval)
}
func (c Config) GetTcpTestTimeOut() time.Duration {
	return time.Duration(c.TcpTestTimeOut) * time.Second
}
