// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareCreateLogic {
	return &ShareCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareCreateLogic) ShareCreate(req *types.ShareCreateRequest, userIdentity string) (resp *types.ShareCreateResponse, err error) {
	if helper.NormalizeInput(req.RepositoryIdentity) == "" {
		return nil, errors.New("资源标识不能为空")
	}
	if !helper.IsValidPositiveDays(req.Expires) {
		return nil, errors.New("分享有效期必须是大于0的天数")
	}

	uuid := helper.GenerateUUID()
	data := &models.ShareLink{
		Identity:           uuid,
		UserIdentity:       userIdentity,
		RepositoryIdentity: req.RepositoryIdentity,
		Expires:            req.Expires,
	}
	_, err = l.svcCtx.Engine.Insert(data)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	_, _ = l.svcCtx.Engine.Where("identity = ?", userIdentity).Get(user)
	repo := new(models.RepositoryPool)
	_, _ = l.svcCtx.Engine.Where("identity = ?", req.RepositoryIdentity).Get(repo)

	fileExt := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(repo.Ext), "."))
	if fileExt == "" {
		fileExt = "none"
	}
	sharerName := strings.TrimSpace(user.Name)
	if sharerName == "" {
		sharerName = userIdentity
	}
	role := user.Role
	if role <= 0 {
		role = 1
	}
	detail := fmt.Sprintf("sharer_name=%s;file_ext=%s;expires_days=%d;share_identity=%s", sharerName, fileExt, req.Expires, uuid)
	helper.AddAuditLog(l.svcCtx.Engine, userIdentity, sharerName, role, "SHARE_CREATE", "share_link", uuid, detail)

	resp = &types.ShareCreateResponse{
		Identity: uuid,
	}
	return
}
