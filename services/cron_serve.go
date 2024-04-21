package services

import (
	"github.com/robfig/cron/v3"
)

type CronFunction struct {
	Function     func(string, string)
	Id           string
	AccountToken string
}

func (s CronFunction) Run() {
	s.Function(s.Id, s.AccountToken)
}

func (s CronFunction) CronServe() {
	cronServer := cron.New(cron.WithSeconds())
	_, err := cronServer.AddJob("@every 30m", s)
	if err != nil {
		return
	}
	cronServer.Start()
	defer cronServer.Stop()
}
