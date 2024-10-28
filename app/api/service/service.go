package service

import (
	"ascale/app/api/conf"
	"ascale/app/api/dao"
	"ascale/pkg/conf/env"
	"ascale/pkg/dlock"
	"ascale/pkg/log"
	"context"
	"runtime"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

// Service struct of service
type Service struct {
	c      *conf.Config
	d      *dao.Dao
	missch chan func()
	dlock  *dlock.Client
	pubsub *pubsub.Client
}

// New create new service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		d:      dao.New(c),
		missch: make(chan func(), 1024*4),
	}
	s.dlock = dlock.New(s.d.Redis())
	var err error
	if env.DeployEnv != env.DeployEnvDev && env.PubsubEndpoint == "" {
		if s.pubsub, err = pubsub.NewClient(context.Background(), env.ProjectID); err != nil {
			log.Fatalf("pubsub.NewClient error(%+v)", err)
		}
	} else {
		options := []option.ClientOption{
			option.WithEndpoint("localhost:8080"),
			option.WithoutAuthentication(),
			option.WithGRPCDialOption(grpc.WithInsecure()),
		}

		if s.pubsub, err = pubsub.NewClient(context.Background(), env.ProjectID, options...); err != nil {
			log.Fatalf("pubsub.NewClient error(%+v)", err)
		}
	}

	s.startSubscriptions()
	s.initialJobCron()
	go s.cacheproc()
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.d.Ping(c)
}

// Close dao.
func (s *Service) Close(ctx context.Context) {
	s.d.Close(ctx)
	s.pubsub.Close()
}

func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			fileName, line := details.FileLine(pc)
			log.Errorf(
				"cacheproc chan full, func(%s), file(%s), line(%d)",
				details.Name(),
				fileName,
				line,
			)
		} else {
			log.Error("cacheproc chan full")
		}
	}
}

func (s *Service) cacheproc() {
	streams := s.fanoutCacheproc(s.missch, 4)

	go s.cacheprocConsumer(streams[0])
	go s.cacheprocConsumer(streams[1])
	go s.cacheprocConsumer(streams[2])
	go s.cacheprocConsumer(streams[3])
}

func (s *Service) cacheprocConsumer(ch <-chan func()) {
	for {
		f := <-ch
		f()
	}
}

func (s *Service) fanoutCacheproc(ch <-chan func(), n int) []chan func() {
	cs := make([]chan func(), n)
	for i := 0; i < n; i++ {
		cs[i] = make(chan func())
	}

	// Distributes the work in a round robin fashion among the stated number
	// of channels until the main channel has been closed. In that case, close
	// all channels and return.
	distributeToChannels := func(ch <-chan func(), cs []chan func()) {
		for {
			for _, c := range cs {
				if val, ok := <-ch; ok {
					c <- val
				} else {
					return
				}
			}
		}
	}

	go distributeToChannels(ch, cs)

	return cs
}
