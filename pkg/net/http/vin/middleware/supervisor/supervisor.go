package supervisor

import (
	"time"

	"ascale/pkg/ecode"
	"ascale/pkg/net/http/vin"
)

// Config supervisor conf.
type Config struct {
	On    bool      // all post/put/delete method off.
	Begin time.Time // begin time
	End   time.Time // end time
}

// Supervisor supervisor midleware.
type Supervisor struct {
	conf *Config
	on   bool
}

// New new and return supervisor midleware.
func New(c *Config) (s *Supervisor) {
	s = &Supervisor{
		conf: c,
	}
	s.Reload(c)
	return
}

// Reload reload supervisor conf.
func (s *Supervisor) Reload(c *Config) {
	if c == nil {
		return
	}
	s.on = c.On && c.Begin.Before(c.End)
	s.conf = c // NOTE datarace but no side effect.
}

func (s *Supervisor) ServeHTTP(c *vin.Context) {
	if s.on {
		now := time.Now()
		method := c.Request.Method
		if s.forbid(method, now) {
			c.JSON(nil, ecode.ServiceUpdate)
			c.Abort()
			return
		}
	}
}

// Handler is router allow handle.
func (s *Supervisor) Handler() vin.HandlerFunc {
	return s.ServeHTTP
}

func (s *Supervisor) forbid(method string, now time.Time) bool {
	// only allow GET request.
	return method != "GET" && now.Before(s.conf.End) && now.After(s.conf.Begin)
}
