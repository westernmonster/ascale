package http

import (
	"ascale/app/api/conf"
	"ascale/app/api/service"
	"ascale/pkg/ecode"
	"ascale/pkg/log"
	"ascale/pkg/net/http/vin"
	"ascale/pkg/net/http/vin/middleware/auth"
)

var (
	srv     *service.Service
	authSvc *auth.Auth
	cnf     *conf.Config
)

func Init(c *conf.Config, s *service.Service) {
	srv = s
	authSvc = auth.New(srv)
	cnf = c

	engine := vin.DefaultServer(c.Vin)
	setupRoute(engine)

	if err := engine.Start(); err != nil {
		log.Fatalf("engine.Start() error(%v)", err)
	}
}

func setupRoute(e *vin.Engine) {
	e.Ping(ping)
	e.Register(register)

	base := e.Group("/")
	route(base)
}

func route(e *vin.RouterGroup) {
	e.GET("/", func(ctx *vin.Context) { ctx.JSON(nil, nil) })
	// e.GET("/system_info", getSystemInfo)
}

// ping check server ok.
func ping(c *vin.Context) {
	var err error
	if err = srv.Ping(c); err != nil {
		log.Errorf("service ping error(%v)", err)
		c.JSON(nil, ecode.ServiceUnavailable)
		return
	}

	c.JSON(nil, nil)
}

// register support discovery.
func register(c *vin.Context) {
	c.JSON(map[string]struct{}{}, nil)
}
