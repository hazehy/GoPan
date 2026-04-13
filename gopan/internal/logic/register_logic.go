// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"log"

	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	req.Name = helper.NormalizeInput(req.Name)
	req.Email = helper.NormalizeInput(req.Email)
	req.Code = helper.NormalizeInput(req.Code)

	if !helper.IsValidUsername(req.Name) {
		return nil, errors.New("用户名格式不正确")
	}
	if !helper.IsValidEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}
	if !helper.IsValidPassword(req.Password) {
		return nil, errors.New("密码长度需在6到32位")
	}

	// 对比验证码
	code, err := l.svcCtx.RDB.Get(l.ctx, req.Email).Result()
	if err != nil {
		return nil, errors.New("验证码已失效")
	}
	if code != req.Code {
		return nil, errors.New("验证码错误")
	}
	// 判断用户是否已注册
	cnt, err := l.svcCtx.Engine.Table("user_basic").Where("name = ?", req.Name).Count(new(models.User))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("用户名已注册")
	}
	// 创建用户
	user := &models.User{
		Identity:           helper.GenerateUUID(),
		Name:               req.Name,
		Password:           helper.Bcrypt(req.Password),
		Email:              req.Email,
		Status:             1,
		Role:               1,
		UploadPermission:   1,
		DownloadPermission: 1,
		SharePermission:    1,
	}
	a, err := l.svcCtx.Engine.Insert(user)
	if err != nil {
		return nil, err
	}
	helper.AddAuditLog(l.svcCtx.Engine, user.Identity, user.Name, user.Role, "USER_REGISTER", "user", user.Identity, "新用户注册")
	log.Println("用户注册成功: ", a)
	return
}
