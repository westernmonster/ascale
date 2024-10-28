package service

import (
	"context"
	"sync"
)

type cronJobFunc func(c context.Context) error

var (
	doOnce   sync.Once
	cronJobs map[string]cronJobFunc = make(map[string]cronJobFunc)
)

func registerCronJob(job string, fn cronJobFunc) {
	cronJobs[job] = fn
}

func (p *Service) initialJobCron() {
	doOnce.Do(func() {
	})
}
