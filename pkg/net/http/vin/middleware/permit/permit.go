package permit

import (
	"ascale/pkg/ecode"
	"ascale/pkg/net/http/vin"
	"ascale/pkg/net/metadata"
)

const (
	_verifyURI             = "/api/session/verify"
	_permissionURI         = "/x/admin/manager/permission"
	_sessIDKey             = "_AJSESSIONID"
	_sessUIDKey            = "uid"      // manager user_id
	_sessUnKey             = "username" // LDAP username
	_defaultDomain         = ".donefirst.com"
	_defaultCookieName     = "ascale-go"
	_defaultCookieLifeTime = 2592000
	// CtxPermissions will be set into ctx.
	CtxPermissions = "permissions"
)

type permissions struct {
	UID   int64    `json:"uid"`
	Perms []string `json:"perms"`
}

type Permit struct {
	sm *SessionManager // user Session
}

type Verify interface {
	Verify() vin.HandlerFunc
}

type Config struct {
	Session *SessionConfig
}

func New() *Permit {
	return &Permit{
		sm: &SessionManager{},
	}
}

func (p *Permit) Verify() vin.HandlerFunc {
	return func(ctx *vin.Context) {
		si, err := p.login(ctx)
		if err != nil {
			ctx.JSON(nil, ecode.Unauthorized)
			ctx.Abort()
			return
		}
		// 存储Session到Redis
		p.sm.SessionRelease(ctx, si)
	}
}

func (p *Permit) login(ctx *vin.Context) (si *Session, err error) {
	si = p.sm.SessionStart(ctx)
	if si.Get(_sessUnKey) == nil {
		var username string
		if username, err = p.verify(ctx); err != nil {
			return
		}
		si.Set(_sessUnKey, username)
	}
	ctx.Set(_sessUnKey, si.Get(_sessUnKey))
	if md, ok := metadata.FromContext(ctx); ok {
		md[metadata.Username] = si.Get(_sessUnKey)
	}
	return
}

func (p *Permit) verify(ctx *vin.Context) (username string, err error) {
	var (
		sid string
		r   = ctx.Request
	)
	session, err := r.Cookie(_sessIDKey)
	if err == nil {
		sid = session.Value
	}
	if sid == "" {
		err = ecode.Unauthorized
		return
	}
	return
}
