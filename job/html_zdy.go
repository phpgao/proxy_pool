package job

import (
	"context"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"math/rand"
	"strings"
	"time"
)

type zdy struct {
	Spider
}

func (s *zdy) Fetch(proxyURL string, useProxy bool) (body string, err error) {
	if s.RandomDelay() {
		time.Sleep(time.Duration(rand.Intn(6)) * time.Second)
	}
	ws, err := util.GetWsFromChrome(fmt.Sprintf("http://%s/json", util.ServerConf.ChromeWS))
	logger.WithField("ws", ws).Debug("get ws addr")
	if err != nil {
		return
	}

	actxt, cancelActxt := chromedp.NewRemoteAllocator(context.Background(), ws)
	defer cancelActxt()

	ctxt, cancelCtxt := chromedp.NewContext(actxt)
	defer cancelCtxt()

	//t, err := chromedp.Targets(ctxt)
	//if err != nil {
	//	logger.WithError(err).Error("targets")
	//}

	if err := chromedp.Run(ctxt,
		network.Enable(),
		network.SetCacheDisabled(true),
		chromedp.Navigate(proxyURL),
		chromedp.WaitVisible("#ipc"),
		chromedp.OuterHTML("html", &body),
		network.ClearBrowserCache(),
		network.ClearBrowserCookies(),
	); err != nil {
		logger.WithError(err).Errorf("Failed getting body of %s", proxyURL)
	}

	body = strings.TrimSpace(body)
	return
}

func (s *zdy) StartUrl() []string {
	return []string{
		"https://www.zdaye.com/FreeIPList.html",
	}
}

func (s *zdy) Enabled() bool {
	return util.ServerConf.ChromeWS != ""
}
func (s *zdy) Cron() string {
	return "@every 5m"
}

func (s *zdy) GetReferer() string {
	return "https://www.zdaye.com"
}

func (s *zdy) Run() {
	getProxy(s)
}

func (s *zdy) Name() string {
	return "zdy"
}

func (s *zdy) Parse(body string) (proxies []*model.HttpProxy, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	list := htmlquery.Find(doc, "//table/tbody/tr[position()>1]")
	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[1]"))
		ip = strings.TrimSpace(ip)
		for _, p := range []string{
			"80", "8080", "8008", "8811", "8888", "9999", "1080", "3000",
		} {
			proxies = append(proxies, &model.HttpProxy{
				Ip:   ip,
				Port: p,
			})
		}

	}
	return
}
