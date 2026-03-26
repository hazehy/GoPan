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

type CodeSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCodeSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CodeSendLogic {
	return &CodeSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CodeSendLogic) CodeSend(req *types.CodeSendRequest) (resp *types.CodeSendResponse, err error) {
	req.Email = helper.NormalizeInput(req.Email)
	if !helper.IsValidEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 判断邮箱是否已注册
	cnt, err := l.svcCtx.Engine.Table("user_basic").Where("email = ?", req.Email).Count(new(models.User))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		err = errors.New("邮箱已注册")
		return nil, err
	}
	// 生成验证码
	code := helper.RandomCode()
	// Redis存储验证码
	l.svcCtx.RDB.Set(l.ctx, req.Email, code, time.Second*time.Duration(define.CodeExpire))
	// 发送验证码
	err = helper.MailCodeSend(req.Email, code)
	if err != nil {
		return nil, err
	}
	return
}
