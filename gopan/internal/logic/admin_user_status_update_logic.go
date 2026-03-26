package logic

import (
	"context"
	"errors"

	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUserStatusUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserStatusUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserStatusUpdateLogic {
	return &AdminUserStatusUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserStatusUpdateLogic) AdminUserStatusUpdate(req *types.AdminUserStatusUpdateRequest) error {
	if req.Status != 1 && req.Status != 2 {
		return errors.New("用户状态仅支持 1(正常) 或 2(禁用)")
	}

	target := new(models.User)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.Identity).Get(target)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("用户不存在")
	}
	if target.Role == 2 {
		return errors.New("不能修改管理员状态")
	}

	_, err = l.svcCtx.Engine.Where("identity = ?", req.Identity).Cols("status").Update(&models.User{Status: req.Status})
	if err == nil {
		helper.AddAuditLog(l.svcCtx.Engine, "SYSTEM", "admin", 2, "USER_STATUS_UPDATE", "user", target.Identity, "管理员更新用户状态")
	}
	return err
}
