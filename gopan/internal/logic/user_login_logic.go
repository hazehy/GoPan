// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"time"

	"gopan/gopan/define"
	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	req.Name = helper.NormalizeInput(req.Name)
	if !helper.IsValidUsername(req.Name) {
		return nil, errors.New("用户名格式不正确")
	}
	if !helper.IsValidPassword(req.Password) {
		return nil, errors.New("密码长度需在6到32位")
	}

	// 查询当前用户
	user := new(models.User)
	has, err := l.svcCtx.Engine.Where("name = ?", req.Name).Get(user)
	if err != nil {
		return nil, err
	}
	if !has || !helper.ComparePassword(user.Password, req.Password) {
		return nil, errors.New("用户名或密码错误")
	}
	if user.Status == 2 {
		return nil, errors.New("当前用户已被禁用")
	}
	// 更新最后一次登录时间
	_, err = l.svcCtx.Engine.Where("id = ?", user.Id).Cols("last_login_at").Update(&models.User{
		LastLoginAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	// 生成token
	token, err := helper.GenerateToken(user.Id, user.Identity, user.Name, user.Role, define.TokenExpire)
	if err != nil {
		return nil, err
	}
	// 生成refresh token
	refreshtoken, err := helper.GenerateToken(user.Id, user.Identity, user.Name, user.Role, define.RefreshTokenExpire)
	if err != nil {
		return nil, err
	}
	resp = new(types.LoginResponse)
	resp.Token = token
	resp.RefreshToken = refreshtoken
	resp.Role = user.Role
	resp.UploadPermission = user.UploadPermission
	resp.DownloadPermission = user.DownloadPermission
	resp.SharePermission = user.SharePermission
	helper.AddAuditLog(l.svcCtx.Engine, user.Identity, user.Name, user.Role, "USER_LOGIN", "user", user.Identity, "用户登录")
	return resp, nil
}
