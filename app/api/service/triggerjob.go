package service

import (
	"ascale/app/api/model"
	"ascale/pkg/def"
	"ascale/pkg/xtime"
	"context"
	"time"
)

func (p *Service) cronSendLittleMessages(c context.Context) (err error) {
	ctx, cancel := context.WithTimeout(c, 1*time.Minute)
	defer cancel()

	// Send a message every 100 millsecond, exit after 1 minutes
	every := 100 * time.Millisecond
	t := time.NewTicker(every)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(every):
			p.Publish(
				context.Background(),
				def.Topics.DoTask,
				&model.DoTaskCommand{Name: "do task"},
			)
		}
	}
}

func (p *Service) TriggerJob(c context.Context, job string) (err error) {
	return p.Publish(
		c,
		def.Topics.Trigger,
		&model.TriggerCommand{Job: job, TriggerTime: xtime.Now().Unix()},
	)
}
