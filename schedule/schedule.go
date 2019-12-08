package schedule

import (
	"github.com/phpgao/proxy_pool/db"
	"github.com/phpgao/proxy_pool/job"
	"github.com/phpgao/proxy_pool/queue"
	"github.com/phpgao/proxy_pool/util"
	"github.com/phpgao/proxy_pool/validator"
	"github.com/robfig/cron/v3"
	"time"
)

type Scheduler struct {
	spiders []job.Crawler
	cronMap map[string]cron.EntryID
	cron    *cron.Cron
}

var (
	config = util.ServerConf
	logger = util.GetLogger("schedule")
)

func (s *Scheduler) Run() {
	logger.Info("adding scheduler...")
	for _, _spider := range s.spiders {
		// trigger once
		go _spider.Run()

		cronID, err := s.cron.AddJob(_spider.Cron(), _spider)

		if err != nil {
			logger.WithError(err).Errorf("error add cron with spider %s", _spider.Name())
		}
		s.cronMap[_spider.Name()] = cronID
	}
	_, _ = s.cron.AddFunc("@every 1m", func() {
		s.report("")
	})
	_, _ = s.cron.AddFunc("@every 1m", validator.Update)
	s.cron.Start()
}

func (s *Scheduler) report(spiderKey string) {
	if spiderKey != "" {
		if entryId, ok := s.cronMap[spiderKey]; ok {
			if ok := s.cron.Entry(entryId).Next.IsZero(); ok {
				logger.Infof("Spider %s hasn't run yet!", spiderKey)
			} else {
				logger.Infof("Next tick of %s --> %s", spiderKey, s.cron.Entry(entryId).Next.Format("2006-01-02 15:04:05"))
			}
		}
	} else {
		for spiderKey, entryId := range s.cronMap {
			if ok := s.cron.Entry(entryId).Next.IsZero(); ok {
				logger.Infof("Spider %s hasn't run yet!", spiderKey)
			} else {
				t := s.cron.Entry(entryId).Next
				if t.Sub(time.Now()) <= 2*time.Minute {
					logger.Infof("Next tick of %s --> %s", spiderKey, s.cron.Entry(entryId).Next.Format("2006-01-02 15:04:05"))
				}
			}
		}
	}

}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		cron:    cron.New(),
		cronMap: make(map[string]cron.EntryID),
	}

	s.spiders = job.GetSpiders(queue.GetNewChan())

	internalJob := Internal{
		channel: queue.OldProxyChan,
		db:      db.GetDb(),
	}

	go internalJob.Run()

	id, err := s.cron.AddJob(config.GetInternalCron(), internalJob)

	if err != nil {
		logger.WithError(err).Error("error initial internal job")
	}
	s.cronMap["internal"] = id
	return s
}
