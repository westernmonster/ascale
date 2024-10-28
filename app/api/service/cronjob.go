package service

import (
	"ascale/app/api/model"
	"ascale/pkg/def"
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
