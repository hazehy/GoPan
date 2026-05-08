package logic

import (
	"context"
	"errors"
	"strings"

	"gopan/gopan/define"
	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type PasswordResetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPasswordResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasswordResetLogic {
	return &PasswordResetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PasswordResetLogic) PasswordReset(req *types.PasswordResetRequest) (resp *types.PasswordResetResponse, err error) {
	req.Email = helper.NormalizeInput(req.Email)
	req.Code = strings.ToLower(helper.NormalizeInput(req.Code))

	if !helper.IsValidEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}
	if !helper.IsValidVerificationCode(req.Code, define.CodeLength) {
		return nil, errors.New("验证码格式不正确")
	}
	if !helper.IsValidPassword(req.Password) {
		return nil, errors.New("密码长度需在6到32位")
	}

	user := new(models.User)
	has, err := l.svcCtx.Engine.Where("email = ? AND deleted_at IS NULL", req.Email).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户不存在")
	}

	resetCodeKey := helper.BuildCodeRedisKey(define.PasswordResetCodeRedisPrefix, req.Email)
	code, err := l.svcCtx.RDB.Get(l.ctx, resetCodeKey).Result()
	if err != nil {
		return nil, errors.New("验证码已失效")
	}
	if code != req.Code {
		return nil, errors.New("验证码错误")
	}

	if _, err = l.svcCtx.Engine.Where("id = ?", user.Id).Cols("password").Update(&models.User{Password: helper.Bcrypt(req.Password)}); err != nil {
		return nil, err
	}
	l.svcCtx.RDB.Del(l.ctx, resetCodeKey)

	helper.AddAuditLog(l.svcCtx.Engine, user.Identity, user.Name, user.Role, "USER_PASSWORD_RESET", "user", user.Identity, "用户通过邮箱验证码重置密码")
	return &types.PasswordResetResponse{}, nil
}
