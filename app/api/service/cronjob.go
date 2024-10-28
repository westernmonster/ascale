package service

import (
	"ascale/app/api/model"
	"ascale/pkg/def"
	"ascale/pkg/xtime"
	"context"
)

func (p *Service) cronDoSmallTask(c context.Context) (err error) {
	p.Publish(
		context.Background(),
		def.Topics.DoSmallTask,
		&model.DoSmallTaskCommand{Name: "small task"},
	)
	return
}

func (p *Service) CronJob(c context.Context, job string) (err error) {
	return p.Publish(
		c,
		def.Topics.CronJob,
		&model.CronJobCommand{Job: job, TriggerTime: xtime.Now().Unix()},
	)
}
