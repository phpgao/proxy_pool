package source

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/apex/log"
	"github.com/parnurzeal/gorequest"
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/util"
	"math/rand"
	"strings"
	"time"
)

var (
	config = util.GetConfig()
	logger = util.GetLogger()
)

func init() {
	htmlquery.DisableSelectorCache = true
}

type Crawler interface {
	Run()
	StartUrl() []string
	Cron() string
	Name() string
	Fetch(string, *model.HttpProxy) (string, error)
	SetProxyChan(chan<- *model.HttpProxy)
	GetProxyChan() chan<- *model.HttpProxy
	Parse(string) ([]*model.HttpProxy, error)
}

type Spider struct {
	ch chan<- *model.HttpProxy
}

func (s *Spider) StartUrl() []string {
	panic("implement me")
}

func (s *Spider) errAndStatus(errs []error, resp gorequest.Response) (err error) {
	if len(errs) > 0 {
		err = errs[0]
		return
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("http code: %d", resp.StatusCode))
	}

	return
}

func (s *Spider) Cron() string {
	panic("implement me")
}
func (s *Spider) TimeOut() int {
	return config.Timeout
}

func (s *Spider) Name() string {
	panic("implement me")
}

func (s *Spider) Parse(string) ([]*model.HttpProxy, error) {
	panic("implement me")
}

func (s *Spider) GetReferer() string {
	return "https://www.baidu.com/"
}

func (s *Spider) SetProxyChan(ch chan<- *model.HttpProxy) {
	s.ch = ch
}

func (s *Spider) GetProxyChan() chan<- *model.HttpProxy {
	return s.ch
}

func (s *Spider) RandomDelay() bool {
	return true
}

func (s *Spider) Fetch(proxyURL string, proxy *model.HttpProxy) (body string, err error) {

	if s.RandomDelay() {
		time.Sleep(time.Duration(rand.Intn(6)) * time.Second)
	}

	request := gorequest.New()
	contentType := "text/html; charset=utf-8"
	var superAgent *gorequest.SuperAgent
	var resp gorequest.Response
	var errs []error
	superAgent = request.Get(proxyURL).
		Set("User-Agent", util.GetRandomUA()).
		Set("Content-Type", contentType).
		Set("Referer", s.GetReferer()).
		Set("Pragma", `no-cache`).
		Timeout(time.Duration(s.TimeOut()) * time.Second)
	if proxy == nil {
		resp, body, errs = superAgent.End()
	} else {
		resp, body, errs = superAgent.Proxy(fmt.Sprintf("http://%s:%s", proxy.Ip, proxy.Port)).End()
	}
	if err = s.errAndStatus(errs, resp); err != nil {
		return
	}

	body = strings.TrimSpace(body)
	return
}

func getProxy(s Crawler) {
	logger.WithField("spider", s.Name()).Debug("spider begin")

	for _, url := range s.StartUrl() {
		go func(proxyURL string, inputChan chan<- *model.HttpProxy) {
			logger.Debugf("Requesting %s", proxyURL)
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()
			var proxyText string
			var newProxies []*model.HttpProxy
			var err error
			fetchFlag := false
			parseFlag := false
			// first try direct
			logger.WithFields(log.Fields{"url": proxyURL}).Debug("try fetching url directly")
			proxyText, err = s.Fetch(proxyURL, nil)
			if err != nil {
				logger.WithFields(log.Fields{"url": proxyURL}).WithError(err).Debug("failed fetching url directly")
			} else {
				fetchFlag = true
			}

			if proxyText != "" {
				newProxies, err = s.Parse(proxyText)
				if err != nil {
					logger.WithFields(log.Fields{"url": proxyURL}).WithError(err).Error("failed parse content by directly")
				}
				if len(newProxies) > 0 {
					parseFlag = true
				}
			}

			if !parseFlag || !fetchFlag {

				r := db.GetDb()
				//if exists proxies,use random max 3 times
				if r.Len() > 0 {
					logger.WithFields(log.Fields{
						"url": proxyURL,
					}).Debug("try fetching url with proxy")
					for i := 1; i <= 3; i++ {
						fetchFlag = false
						parseFlag = false
						randomProxy, err := r.Random()
						if err != nil {
							logger.WithFields(log.Fields{
								"url":  proxyURL,
								"time": i,
							}).WithError(err).Error("failed get random proxy")
							continue
						}
						logger.WithFields(log.Fields{
							"url":   proxyURL,
							"time":  i,
							"proxy": fmt.Sprintf("%s:%s", randomProxy.Ip, randomProxy.Port),
						}).Debug("try fetching url with proxy")

						proxyText, err = s.Fetch(proxyURL, &randomProxy)
						if err != nil {
							logger.WithFields(log.Fields{
								"url":      proxyURL,
								"try_time": i,
								"proxy":    fmt.Sprintf("%s:%s", randomProxy.Ip, randomProxy.Port),
							}).WithError(err).Debug("failed fetching url with proxy")
							continue
						}
						fetchFlag = true
						if proxyText != "" {
							newProxies, err = s.Parse(proxyText)
							if err != nil {
								logger.WithFields(log.Fields{
									"url": proxyURL,
								}).WithError(err).Error("failed parse content by proxy")
							}
							if len(newProxies) > 0 {
								logger.WithFields(log.Fields{
									"url":   proxyURL,
									"time":  i,
									"proxy": fmt.Sprintf("%s:%s", randomProxy.Ip, randomProxy.Port),
								}).Debug("success fetching url with proxy")
								parseFlag = true
								break
							}
						}

					}
				} else {
					logger.WithFields(log.Fields{
						"url": proxyURL,
					}).Debug("no proxy available,pass")
				}
			}

			// if empty content or failed
			// user proxy 3 times
			var tmpMap = map[string]int{}
			if parseFlag && fetchFlag {
				for _, newProxy := range newProxies {
					newProxy.Ip = strings.TrimSpace(newProxy.Ip)
					newProxy.Port = strings.TrimSpace(newProxy.Port)
					if _, found := tmpMap[newProxy.GetKey()]; found {
						continue
					}
					tmpMap[newProxy.GetKey()] = 1
					newProxy.From = s.Name()
					if newProxy.Score == 0 {
						newProxy.Score = config.Score
					}
					if util.FilterProxy(newProxy) {
						inputChan <- newProxy
					}
				}
			}

		}(url, s.GetProxyChan())
	}

	logger.WithField("spider", s.Name()).Debug("send jobs done")
}
