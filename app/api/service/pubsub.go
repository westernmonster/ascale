package service

import (
	"ascale/pkg/def"
	"ascale/pkg/log"
	"ascale/pkg/stat/prom"
	"ascale/pkg/xtime"
	"context"
	"fmt"
	"reflect"
	"time"

	"cloud.google.com/go/pubsub"
	jsoniter "github.com/json-iterator/go"
)

func (p *Service) EnsureTopic(ctx context.Context, topic string) (*pubsub.Topic, error) {
	ret := p.pubsub.Topic(topic)
	exists, err := ret.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return p.pubsub.CreateTopic(ctx, topic)
	}
	return ret, nil
}

func (p *Service) EnsureSubscription(
	ctx context.Context,
	topic string,
	deadleterPolicy *pubsub.DeadLetterPolicy,
) (*pubsub.Subscription, error) {
	// hostname := env.GetHostname()

	subID := fmt.Sprintf("%s.sub.%s", topic, "ascale")
	ret := p.pubsub.Subscription(subID)
	exists, err := ret.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		cfg := pubsub.SubscriptionConfig{
			Topic: p.pubsub.Topic(topic),
		}
		if deadleterPolicy != nil {
			cfg.DeadLetterPolicy = deadleterPolicy
		}
		return p.pubsub.CreateSubscription(ctx, subID, cfg)
	}
	return ret, nil
}

// CleanupTopic deletes a topic with all subscriptions and logs any
// error. Useful for defer.
func (p *Service) CleanupTopic(ctx context.Context, project, topic string) {
	if err := p.pubsub.Topic(topic).Delete(ctx); err != nil {
		log.For(ctx).Errorf("Failed to delete topic %v: %v", topic, err)
	}
}

// Publish is a simple utility for publishing a set of string messages
// serially to a pubsub topic. Small scale use only.
func (p *Service) Publish(ctx context.Context, topic string, msg interface{}) (err error) {
	// t, err := p.EnsureTopic(ctx, topic)
	// if err != nil {
	// 	log.For(ctx).Errorf("Publish() error(%+v)", err)
	// 	return err
	// }
	var data []byte
	if data, err = jsoniter.Marshal(msg); err != nil {
		log.For(ctx).Errorf("Publish() error(%+v)", err)
		return
	}

	_, err = p.pubsub.Topic(topic).Publish(ctx, &pubsub.Message{
		Data: data,
	}).Get(ctx)

	if err != nil {
		log.For(ctx).Errorf("Publish() topic(%s) msg(%+v) error(%+v)", topic, msg, err)
		return
	}
	return
}

func (p *Service) getAllTopics() []string {
	t := reflect.ValueOf(def.Topics)
	names := make([]string, t.NumField())
	for i := range names {
		names[i] = t.Field(i).String()
	}

	return names
}

func (p *Service) subscriptions() {
	ctx := context.Background()
	var err error

	topics := p.getAllTopics()

	for _, v := range topics {
		if _, err = p.EnsureTopic(ctx, v); err != nil {
			log.Fatalf("subscription topic(%s)  error(%+v)", v, err)
		}
	}

	createSubscription := func(c context.Context, topic string, maxOutstandingMessages int, deadPolicy *pubsub.DeadLetterPolicy, task func(ctx context.Context, msg *pubsub.Message)) {
		go func() {
			sub, err := p.EnsureSubscription(c, topic, deadPolicy)
			if err != nil {
				log.Fatalf("subscription topic(%s)  error(%+v)", topic, err)
			}

			sub.ReceiveSettings.Synchronous = true
			sub.ReceiveSettings.MaxOutstandingMessages = maxOutstandingMessages
			sub.ReceiveSettings.MaxOutstandingBytes = 1e10

			i := 0
			for {
				time.Sleep(time.Second * 2)
				if err = sub.Receive(c, task); err != nil {
					log.For(c).Errorf("subscription topic Receive (%s) error(%+v)", topic, err)
					i++
					if i > 10 {
						log.Fatalf("subscription topic(%s) Receive  error(%+v)", topic, err)
					}
				}
			}
		}()
	}

	// deadPolicy := &pubsub.DeadLetterPolicy{
	// 	DeadLetterTopic: fmt.Sprintf(
	// 		"projects/%s/topics/%s",
	// 		env.ProjectID,
	// 		def.Topics.DeadLetter,
	// 	),
	// 	MaxDeliveryAttempts: 5,
	// }

	createSubscription(ctx, def.Topics.Trigger, 1, nil, p.jobTrigger)

	createSubscription(ctx, def.Topics.DoTask, 1, nil, p.jobDoTask)

	// DeadLetter
	createSubscription(ctx, def.Topics.DeadLetter, 1, nil, p.logDeadLetter)
}

func (p *Service) logDeadLetter(c context.Context, msg *pubsub.Message) {
	now := xtime.Now()

	defer func() {
		prom.Consumer.Timing(
			fmt.Sprintf("consumer:%s", def.Topics.DeadLetter),
			int64(time.Since(now)/time.Millisecond),
		)
		prom.Consumer.Incr(fmt.Sprintf("consumer:%s", def.Topics.DeadLetter))
	}()

	log.For(c).Errorf("DeadLetter, data(%s)", string(msg.Data))
	msg.Ack()
}
