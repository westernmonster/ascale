package permit

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"sync"
	"time"

	"ascale/pkg/cache/redis"
	"ascale/pkg/log"
	"ascale/pkg/net/http/vin"

	jsoniter "github.com/json-iterator/go"
)

type Session struct {
	Sid string

	lock   sync.RWMutex
	Values map[string]interface{}
}

type SessionConfig struct {
	SessionIDLength int
	CookieLifeTime  int
	CookieName      string
	Domain          string

	Redis *redis.Config
}

type SessionManager struct {
	redis *redis.Pool // Session cache
	c     *SessionConfig
}

func newSessionManager(c *SessionConfig) (s *SessionManager) {
	s = &SessionManager{
		redis: redis.NewPool(c.Redis),
		c:     c,
	}
	return
}

func (s *SessionManager) SessionStart(ctx *vin.Context) (si *Session) {
	// check manager Session id, if err or no exist need new one.
	if si, _ = s.cache(ctx); si == nil {
		si = s.newSession(ctx)
	}
	return
}

// SessionRelease flush session into store.
func (s *SessionManager) SessionRelease(ctx *vin.Context, sv *Session) {
	// set http cookie
	s.setHTTPCookie(ctx, s.c.CookieName, sv.Sid)
	// set mc
	conn := s.redis.Get()
	defer conn.Close()
	key := sv.Sid

	if v, err := jsoniter.Marshal(sv); err != nil {
		log.For(ctx).Errorf("SessionManager set error(%s,%v)", key, err)
		return
	} else {
		if err := conn.Send("SET", key, v); err != nil {
			log.For(ctx).Errorf("SessionManager set error(%s,%v)", key, err)
			return
		}

		if err := conn.Send("EXPIRE", key, int32(s.c.CookieLifeTime)); err != nil {
			log.For(ctx).Errorf("SessionManager set error(%s,%v)", key, err)
			return
		}
		if err := conn.Flush(); err != nil {
			log.For(ctx).Errorf("SessionManager set error(%v)", err)
			return
		}
	}
}

// SessionDestroy destroy session.
func (s *SessionManager) SessionDestroy(ctx *vin.Context, sv *Session) {
	conn := s.redis.Get()
	defer conn.Close()
	if err := conn.Send("DEL", sv.Sid); err != nil {
		log.For(ctx).Errorf("SessionManager delete error(%s,%v)", sv.Sid, err)
	}
}

func (s *SessionManager) cache(ctx *vin.Context) (res *Session, err error) {
	ck, err := ctx.Request.Cookie(s.c.CookieName)
	if err != nil || ck == nil {
		return
	}
	sid := ck.Value
	// get from cache
	conn := s.redis.Get()
	defer conn.Close()

	var data []byte
	if data, err = redis.Bytes(conn.Do("GET", sid)); err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		} else {
			log.For(ctx).Errorf("conn.Do(GET, %s) error(%v)", sid, err)
			return
		}
	}

	res = &Session{}
	if err = jsoniter.Unmarshal(data, res); err != nil {
		log.For(ctx).Errorf("jsoniter.Unmarshal%v) error(%v)", string(data), err)
	}
	return
}

func (s *SessionManager) newSession(ctx context.Context) (res *Session) {
	b := make([]byte, s.c.SessionIDLength)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return nil
	}
	res = &Session{
		Sid:    hex.EncodeToString(b),
		Values: make(map[string]interface{}),
	}
	return
}

func (s *SessionManager) setHTTPCookie(ctx *vin.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     "/",
		HttpOnly: true,
		Domain:   _defaultDomain,
	}
	cookie.MaxAge = _defaultCookieLifeTime
	cookie.Expires = time.Now().Add(time.Duration(_defaultCookieLifeTime) * time.Second)
	http.SetCookie(ctx.Writer, cookie)
}

// Get get value by key.
func (s *Session) Get(key string) (value interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	value = s.Values[key]
	return
}

// Set set value into session.
func (s *Session) Set(key string, value interface{}) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Values[key] = value
	return
}
