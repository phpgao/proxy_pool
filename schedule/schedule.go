package schedule

import "C"
import (
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/model"
	"github.com/phpgao/proxy_pool/spider"
	"github.com/phpgao/proxy_pool/util"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	spiders []spider.Crawler
	cronMap map[string]cron.EntryID
	cron    *cron.Cron
}

var (
	config = util.GetConfig()
	logger = util.GetLogger()
)

func init() {
	logger.Info("init scheduler...")
}

func (s *Scheduler) Run() {
	logger.Info("adding scheduler...")
	for _, _spider := range s.spiders {
		// trigger once
		go _spider.Run()

		cronID, err := s.cron.AddJob(_spider.Cron(), _spider)

		if err != nil {
			logger.WithError(err).Errorf("Error when add cron with spider %s", _spider.Name())
		}
		s.cronMap[_spider.Name()] = cronID
	}
	_, _ = s.cron.AddFunc("@every 1m", func() {
		s.report("")
	})
	s.cron.Start()
	s.report("")
}

func (s *Scheduler) report(spiderKey string) {
	if spiderKey != "" {
		entryId := s.cronMap[spiderKey]
		if ok := s.cron.Entry(entryId).Next.IsZero(); ok {
			logger.Infof("Spider % hasn't run yet!")
		} else {
			logger.Infof("Next tick of %s --> %s", spiderKey, s.cron.Entry(entryId).Next.Format("2006-01-02 15:04:05"))
		}
	} else {
		for spiderKey, entryId := range s.cronMap {
			if ok := s.cron.Entry(entryId).Next.IsZero(); ok {
				logger.Infof("Spider % hasn't run yet!")
			} else {
				logger.Infof("Next tick of %s --> %s", spiderKey, s.cron.Entry(entryId).Next.Format("2006-01-02 15:04:05"))
			}
		}
	}

}

// write only chan
func NewScheduler(newChan chan *model.HttpProxy, oldChan chan *model.HttpProxy, db db.Store) *Scheduler {
	s := &Scheduler{
		cron:    cron.New(),
		cronMap: make(map[string]cron.EntryID),
	}

	s.spiders = spider.GetSpiders(newChan)

	internalJob := Internal{
		channel: oldChan,
		db:      db,
	}
	go internalJob.Run()

	id, err := s.cron.AddJob(config.GetInternalCron(), internalJob)

	if err != nil {
		logger.WithError(err).Error("error when init internal job")
	}
	s.cronMap["internal"] = id
	return s
}
