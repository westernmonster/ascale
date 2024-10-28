package antispam_test

import (
	"time"

	"ascale/pkg/cache/redis"
	"ascale/pkg/net/http/vin"
	"ascale/pkg/net/http/vin/middleware/antispam"
	"ascale/pkg/xtime"
)

// This example create a antispam middleware instance and attach to a vin engine,
// it will protect '/ping' API with specified policy.
// If anyone who requests this API more frequently than 1 req/second or 1 req/hour,
// a StatusServiceUnavailable error will be raised.
func Example() {
	anti := antispam.New(&antispam.Config{
		On:     true,
		Second: 1,
		N:      1,
		Hour:   1,
		M:      1,
		Redis: &redis.Config{
			MaxActive:    10,
			MaxIdle:      10,
			IdleTimeout:  xtime.Duration(time.Second * 60),
			Name:         "test",
			Proto:        "tcp",
			Addr:         "172.18.33.60:6889",
			DialTimeout:  xtime.Duration(time.Second),
			ReadTimeout:  xtime.Duration(time.Second),
			WriteTimeout: xtime.Duration(time.Second),
		},
	})

	engine := vin.Default()
	engine.Use(anti.Handler())
	engine.GET("/ping", func(c *vin.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":18080")
}
