package supervisor_test

import (
	"time"

	"ascale/pkg/net/http/vin"
	"ascale/pkg/net/http/vin/middleware/supervisor"
)

// This example create a supervisor middleware instance and attach to a mars engine,
// it will allow '/ping' API can be requested with specified policy.
// This example will block all http method except `GET` on '/ping' API in next hour,
// and allow in further.
func Example() {
	now := time.Now()
	end := now.Add(time.Hour * 1)
	spv := supervisor.New(&supervisor.Config{
		On:    true,
		Begin: now,
		End:   end,
	})

	engine := vin.Default()
	engine.Use(spv.Handler())
	engine.GET("/ping", func(c *vin.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":18080")
}
