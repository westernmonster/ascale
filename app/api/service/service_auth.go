package service

import (
	"ascale/pkg/net/http/vin/middleware/auth"
	"context"
)

func (p *Service) GetTokenInfo(c context.Context, token string) (r *auth.AuthReply, err error) {
	return
}

func (p *Service) UpdateTokenInfo(c context.Context, token string) (err error) {
	return
}
