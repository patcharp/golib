package crontab

import (
	"github.com/robfig/cron/v3"
)

type Service struct {
	c *cron.Cron
}

func NewService() Service {
	return Service{c: cron.New()}
}

func (s *Service) Start() {
	s.c.Start()
}

func (s *Service) Stop() {
	if s.c != nil {
		s.c.Stop()
	}
}

func (s *Service) AddJob(spec string, cmd func()) (cron.EntryID, error) {
	return s.c.AddFunc(spec, cmd)
}

func (s *Service) RemoveJob(entryId cron.EntryID) {
	s.c.Remove(entryId)
}

func (s *Service) GetEntryJobs() []cron.Entry {
	return s.c.Entries()
}
