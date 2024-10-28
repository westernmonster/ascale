// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vin

import (
	"ascale/pkg/def"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"ascale/pkg/ecode"
	"ascale/pkg/log"
	"ascale/pkg/net/metadata"
)

func Logger() HandlerFunc {
	const noUser = "no_user"
	return func(c *Context) {
		now := time.Now()
		ip := metadata.String(c, metadata.RemoteIP)
		req := c.Request
		path := req.URL.Path
		params := req.Form
		headers := c.Request.Header

		var quota float64
		if deadline, ok := c.Context.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		c.Next()

		uid, ok := metadata.Value(c, metadata.Uid).(int64)
		if !ok {
			uid = 0
		}

		var err error

		lastErr := c.Errors.Last()
		if lastErr != nil && lastErr.Err != nil {
			err = lastErr
		}

		cerr := ecode.Cause(err)
		dt := time.Since(now)
		caller := metadata.String(c, metadata.Caller)
		if caller == "" {
			caller = noUser
		}

		stats.Incr(caller, path[1:], strconv.FormatInt(int64(cerr.Code()), 10))
		stats.Timing(caller, int64(dt/time.Millisecond), path[1:])

		errmsg := ""
		isSlow := dt >= (time.Millisecond * 1000)

		appPlatform := headers.Get(def.AppHeader.Platform)
		appVersion := headers.Get(def.AppHeader.Version)
		fields := []zapcore.Field{
			zap.String("method", req.Method),
			zap.String("uid", strconv.FormatInt(uid, 10)),
			zap.String("ip", ip),
			zap.String("user", caller),
			zap.String("path", path),
			zap.String("params", params.Encode()),
			zap.Int("ret", cerr.Code()),
			zap.String("msg", cerr.Message()),
			zap.String("stack", fmt.Sprintf("%+v", err)),
			zap.Float64("timeout_quota", quota),
			zap.Float64("ts", dt.Seconds()),
			zap.String("source", "http-access-log"),
			zap.String("appPlatform", appPlatform),
			zap.String("appVersion", appVersion),
		}

		if err != nil {
			errmsg = err.Error()
			fields = append(fields, zap.String("err", errmsg))
			if cerr.Code() > 0 {
				log.For(c).WarnWithFields("http", fields...)
				return
			}
			log.For(c).ErrorWithFields("http", fields...)
			return

		} else {
			if isSlow {
				log.For(c).WarnWithFields("http", fields...)
				return
			}
		}

		if _ignorePaths[path] {
			return
		}
		log.For(c).InfoWithFields("http", fields...)
	}
}
