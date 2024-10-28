package vin

import (
	"ascale/pkg/def"
	"ascale/pkg/xtime"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func UUID() HandlerFunc {
	return func(c *Context) {
		ck, err := c.Request.Cookie(def.DoneUUID)
		if err == nil && ck.Value != "" {
			c.Next()
		} else {
			cookie := new(http.Cookie)
			cookie.Name = def.DoneUUID
			cookie.MaxAge = 0
			cookie.HttpOnly = false
			cookie.Path = "/"
			cookie.Expires = xtime.Now().AddDate(50, 0, 0)
			cookie.Value = uuid.NewV4().String()
			http.SetCookie(c.Writer, cookie)
			c.Next()
		}
	}
}
