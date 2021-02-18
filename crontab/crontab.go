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

func (s *Service) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return s.c.AddFunc(spec, cmd)
}

func (s *Service) AddJob(spec string, cmd cron.Job) (cron.EntryID, error) {
	return s.c.AddJob(spec, cmd)
}

func (s *Service) RemoveJob(entryId int) {
	s.c.Remove(cron.EntryID(entryId))
}

func (s *Service) GetEntryJob(entryId int) cron.Entry {
	return s.c.Entry(cron.EntryID(entryId))
}

func (s *Service) GetEntryJobs() []cron.Entry {
	return s.c.Entries()
}
