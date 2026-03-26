// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"gopan/gopan/define"
	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokenRefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTokenRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenRefreshLogic {
	return &TokenRefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokenRefreshLogic) TokenRefresh(req *types.LoginRequest, authorization string) (resp *types.LoginResponse, err error) {
	uc, err := helper.AnalyzeToken(authorization)
	if err != nil {
		return nil, err
	}
	token, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, uc.Role, define.TokenExpire)
	if err != nil {
		return nil, err
	}
	refreshtoken, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, uc.Role, define.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}
	resp = new(types.LoginResponse)
	resp.Token = token
	resp.RefreshToken = refreshtoken
	resp.Role = uc.Role
	return
}
