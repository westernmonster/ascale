package auth

import (
	"ascale/pkg/database/sqlx/types"
	"ascale/pkg/def"
	"ascale/pkg/ecode"
	"ascale/pkg/net/http/vin"
	"ascale/pkg/net/metadata"
	"context"
	"net/http"
	"strings"
)

type IAuth interface {
	GetTokenInfo(ctx context.Context, token string) (reply *AuthReply, err error)
	UpdateTokenInfo(ctx context.Context, token string) (err error)
}

type AuthReply struct {
	Login             bool          `json:"login"`
	Role              string        `json:"role"`
	Uid               int64         `json:"uid,string"`
	Expires           int64         `json:"expires"`
	Impersonate       bool          `json:"impersonate"`
	EnableAdminAccess types.BitBool `json:"enable_admin_access"`
}

type Auth struct {
	Identity IAuth
}

type authFunc func(*vin.Context) (int64, error)

func New(au IAuth) *Auth {
	auth := &Auth{
		Identity: au,
	}
	return auth
}

func (a *Auth) User(c *vin.Context) {
	req := c.Request
	cookie, _ := req.Cookie("token")
	if cookie != nil {
		a.UserWeb(c)
	} else {
		a.UserMobile(c)
	}
}

// UserWeb is used to mark path as web access required.
func (a *Auth) UserWeb(ctx *vin.Context) {
	a.midAuth(ctx, a.AuthCookie, def.Origin.Web)
}

// UserMobile is used to mark path as mobile access required.
func (a *Auth) UserMobile(ctx *vin.Context) {
	a.midAuth(ctx, a.AuthToken, def.Origin.App)
}

// AuthToken is used to authorize request by token
func (a *Auth) AuthToken(ctx *vin.Context) (accountID int64, err error) {
	tokenStr := ""
	if v, ok := getBearer(ctx.GetHeader("Authorization")); ok {
		tokenStr = v
	}

	if tokenStr == "" {
		return 0, ecode.NoLogin
	}

	reply, err := a.Identity.GetTokenInfo(ctx, tokenStr)
	if err != nil {
		return 0, err
	}

	if !reply.Login {
		return 0, ecode.NoLogin
	}

	if reply.Role != def.UserRole.User {
		return 0, ecode.NoLogin
	}

	if reply.Impersonate && ctx.Request.Method != http.MethodGet {
		return 0, ecode.MethodNoPermission
	}

	return reply.Uid, nil
}

// AuthCookie is used to authorize request by cookie
func (a *Auth) AuthCookie(ctx *vin.Context) (int64, error) {
	cookie, err := ctx.Request.Cookie("token")
	if err != nil {
		return 0, ecode.NoLogin
	}

	reply, err := a.Identity.GetTokenInfo(ctx, cookie.Value)
	if err != nil {
		return 0, err
	}

	if !reply.Login {
		return 0, ecode.NoLogin
	}

	if reply.Role != def.UserRole.User {
		return 0, ecode.NoLogin
	}

	if reply.Impersonate && ctx.Request.Method != http.MethodGet {
		return 0, ecode.MethodNoPermission
	}

	return reply.Uid, nil
}

func (a *Auth) midAuth(ctx *vin.Context, auth authFunc, origin string) {
	accountID, err := auth(ctx)
	if err != nil {
		if err == ecode.MethodNoPermission {
			ctx.JSON(nil, err)
		} else {
			ctx.JSON(nil, ecode.Unauthorized)
		}
		ctx.Abort()
		return
	}
	setAccountID(ctx, accountID, origin)
}

// set account id into context
func setAccountID(ctx *vin.Context, id int64, origin string) {
	ctx.Set("uid", id)
	ctx.Set("origin", origin)
	if md, ok := metadata.FromContext(ctx); ok {
		md[metadata.Uid] = id
		md[metadata.Origin] = origin
		return
	}
}

func getBearer(auth string) (jwt string, ok bool) {
	ret := strings.Split(auth, " ")
	if len(ret) == 2 && ret[0] == "Bearer" {
		return ret[1], true
	}
	return "", false
}
