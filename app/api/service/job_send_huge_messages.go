package service

import (
	"ascale/app/api/model"
	"ascale/pkg/def"
	"context"
	"fmt"
	"time"
)

func (p *Service) triggerSendHugeAmountMessages(c context.Context) (err error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Minute)
	defer cancel()

	// Send a message every 10 millsecond, exit after 10 minutes
	every := 10 * time.Millisecond
	t := time.NewTicker(every)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(every):
			fmt.Println("-=====================")
			p.Publish(
				context.Background(),
				def.Topics.DoTask,
				&model.DoTaskCommand{Name: "do task"},
			)
		}
	}
}
