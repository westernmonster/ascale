package rate_test

import (
	"ascale/pkg/net/http/vin"
	"ascale/pkg/net/http/vin/middleware/rate"
)

// This example create a rate middleware instance and attach to a vin engine,
// it will protect '/ping' API frequency with specified policy.
// If any internal service who requests this API more frequently than 1 req/second,
// a StatusTooManyRequests error will be raised.
func Example() {
	lim := rate.New(&rate.Config{
		URLs: map[string]*rate.Limit{
			"/ping": &rate.Limit{Limit: 1, Burst: 2},
		},
		Apps: map[string]*rate.Limit{
			"a-secret-app-key": &rate.Limit{Limit: 1, Burst: 2},
		},
	})

	engine := vin.Default()
	engine.Use(lim.Handler())
	engine.GET("/ping", func(c *vin.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":18080")
}
