package service

import (
	"ascale/app/api/model"
	"ascale/pkg/def"
	"ascale/pkg/dlock"
	"ascale/pkg/log"
	"ascale/pkg/stat/prom"
	"ascale/pkg/xtime"
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	jsoniter "github.com/json-iterator/go"
)

type cronJobFunc func(c context.Context) error

var (
	doOnce   sync.Once
	cronJobs map[string]cronJobFunc = make(map[string]cronJobFunc)
)

func registerTriggerJob(job string, fn cronJobFunc) {
	cronJobs[job] = fn
}

func (p *Service) initialTriggerJob() {
	doOnce.Do(func() {
		registerTriggerJob(def.TriggerJob.CronSendLittleMessage, p.cronSendLittleMessages)
		registerTriggerJob(def.TriggerJob.SendHugeMessage, p.triggerSendHugeAmountMessages)
	})
}

func (p *Service) doPublishConcurrently(
	ctx context.Context,
	ids []int64,
	concurrency int,
	fn func(v int64),
) {
	var wg sync.WaitGroup
	wg.Add(concurrency)
	ch := make(chan int64)

	// goroutine
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for v := range ch {
				fn(v)
			}
		}()
	}

	// Send message to channel
	go func() {
		for _, v := range ids {
			ch <- v
		}
		close(ch)
	}()

	// Wait for all goroutines to complete
	wg.Wait()
}

func (p *Service) jobTrigger(c context.Context, msg *pubsub.Message) {
	var err error
	cmd := new(model.TriggerCommand)
	if err = jsoniter.Unmarshal(msg.Data, cmd); err != nil {
		log.For(c).Errorf("jobCron error(%+v)", err)
		msg.Ack()
		return
	}
	log.For(c).Infof("jobCron.start job(%+v), trigger(%d)", cmd.Job, cmd.TriggerTime)
	fn, ok := cronJobs[cmd.Job]
	if !ok {
		log.For(c).Errorf("p.jobCron(%s) not found.", cmd.Job)
		msg.Ack()
		return
	}
	msg.Ack()

	now := time.Now()
	beginFun := xtime.NowUnix()
	defer func() {
		prom.TriggerJob.Timing(
			fmt.Sprintf("cronjob:%s", cmd.Job),
			int64(time.Since(now)/time.Millisecond),
		)
		prom.TriggerJob.Incr(fmt.Sprintf("cronjob:%s", cmd.Job))
		endFun := xtime.NowUnix()
		if endFun-beginFun > 60 {
			log.For(c).
				Warnf("p.jobCron.duration job (%+v) start(%+v) too long(%+v)", cmd.Job, now, endFun-beginFun)
		}
	}()

	var lock *dlock.Lock
	if lock, err = p.dlock.Obtain(c, def.CronJobLock(cmd.Job), 30*time.Second, &dlock.Options{Context: c}); err != nil {
		log.For(c).Errorf("obtain jobCron  lock failed,job (%+v) error(%+v) ", cmd.Job, err)
		err = nil
		return
	}
	defer lock.Release(c)

	if err = fn(c); err != nil {
		log.For(c).Errorf("p.jobCron(%s) error(%+v)", cmd.Job, err)
		msg.Ack()
		return
	}

	log.For(c).Infof("jobCron.Ack job(%+v), trigger(%d)", cmd.Job, cmd.TriggerTime)

	msg.Ack()

	return
}
