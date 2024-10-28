package service

import (
	"ascale/app/api/model"
	"ascale/pkg/def"
	"ascale/pkg/log"
	"ascale/pkg/stat/prom"
	"ascale/pkg/xtime"
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	jsoniter "github.com/json-iterator/go"
)

func (p *Service) jobDoTask(c context.Context, msg *pubsub.Message) {
	now := xtime.Now()
	var err error
	defer func() {
		msg.Ack()
		prom.Consumer.Timing(
			fmt.Sprintf("consumer:%s", def.Topics.DoTask),
			int64(time.Since(now)/time.Millisecond),
		)
		prom.Consumer.Incr(fmt.Sprintf("consumer:%s", def.Topics.DoTask))
	}()

	cmd := new(model.DoTaskCommand)
	if err = jsoniter.Unmarshal(msg.Data, cmd); err != nil {
		log.For(c).Errorf("jobSendMail error(%+v)", err)
		return
	}

	log.For(c).Info("Do some small task")
}
