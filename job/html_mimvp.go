package job

import (
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"os/exec"
)

func init() {
	spider := mimvpProxy{}
	if spider.Enabled() {
		register(spider)
	} else {
		util.GetLogger("spider").Info("spider mimvpProxy disabled")
	}
}

type mimvpProxy struct {
	Spider
}

func (s mimvpProxy) Cron() string {
	return "@every 2m"
}

func (s mimvpProxy) GetReferer() string {
	return "https://proxy.mimvp.com/"
}

func (s mimvpProxy) StartUrl() []string {
	return []string{
		"https://proxy.mimvp.com/freesecret",
		"https://proxy.mimvp.com/freesecret?proxy=in_hp&sort=&page=2",
		"https://proxy.mimvp.com/freesecret?proxy=in_hp&sort=&page=3",
		"https://proxy.mimvp.com/freesecret?proxy=in_hp&sort=&page=4",
		"https://proxy.mimvp.com/freesecret?proxy=in_hp&sort=&page=5",
	}
}

func (s mimvpProxy) Enabled() bool {
	_, err := exec.LookPath("tesseract")
	if err != nil {
		return false
	}
	return true
}

func (s mimvpProxy) Run() {
	getProxy(&s)
}

func (s mimvpProxy) Name() string {
	return "mimvpProxy"
}

func (s mimvpProxy) Parse(body string) (proxies []*model.HttpProxy, err error) {

	return
}
