package main

import (
	"ascale/app/api/conf"
	"ascale/app/api/http"
	"ascale/app/api/service"
	ecode "ascale/pkg/ecode/tip"
	"ascale/pkg/gid"
	"ascale/pkg/log"
	"ascale/pkg/tracing"
	"context"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"
)

var svc *service.Service

func main() {
	flag.Parse()
	// Load Config
	if err := conf.Init(); err != nil {
		log.Fatalf("conf.Init() error(%v)", err)
	}

	// Init Global ID
	if err := gid.Init(); err != nil {
		log.Fatalf("gid.Init() error(%v)", err)
	}

	// init ecode
	ecode.Init()

	// init log
	log.Init(conf.Conf.Log)
	defer log.Close()

	log.Info("ascale-api start")
	// init trace
	tracing.Init(conf.Conf.Tracer)

	// init service
	svc = service.New(conf.Conf)
	http.Init(conf.Conf, svc)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("ascale-api get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close(context.Background())
			log.Info("ascale-api exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
