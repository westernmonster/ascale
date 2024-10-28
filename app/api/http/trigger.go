package http

import (
	"ascale/app/api/model"
	"ascale/pkg/ecode"
	"ascale/pkg/log"
	"ascale/pkg/net/http/vin"
)

func triggerJob(c *vin.Context) {
	arg := new(model.ArgJob)
	if e := c.BindJSON(arg); e != nil {
		return
	}

	if e := arg.Validate(); e != nil {
		log.For(c).Warnf("arg.Validate() error(%+v)", e)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.TriggerJob(c, arg.Job))
}
