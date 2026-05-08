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

type PasswordResetCodeSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPasswordResetCodeSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasswordResetCodeSendLogic {
	return &PasswordResetCodeSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PasswordResetCodeSendLogic) PasswordResetCodeSend(req *types.PasswordResetCodeSendRequest) (resp *types.PasswordResetCodeSendResponse, err error) {
	req.Email = helper.NormalizeInput(req.Email)
	if !helper.IsValidEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	cnt, err := l.svcCtx.Engine.Table("user_basic").Where("email = ?", req.Email).Count(new(models.User))
	if err != nil {
		return nil, err
	}
	if cnt == 0 {
		return nil, errors.New("邮箱未注册")
	}

	code := helper.RandomCode()
	resetCodeKey := helper.BuildCodeRedisKey(define.PasswordResetCodeRedisPrefix, req.Email)
	if err = l.svcCtx.RDB.Set(l.ctx, resetCodeKey, code, time.Second*time.Duration(define.CodeExpire)).Err(); err != nil {
		return nil, err
	}
	if err = helper.MailCodeSend(req.Email, code); err != nil {
		return nil, err
	}

	return &types.PasswordResetCodeSendResponse{}, nil
}
