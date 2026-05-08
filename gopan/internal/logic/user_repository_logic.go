// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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

type UserRepositoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRepositoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRepositoryLogic {
	return &UserRepositoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRepositoryLogic) UserRepository(req *types.UserRepositoryRequest, userIdentity string) (resp *types.UserRepositoryResponse, err error) {
	if err := ensureUserCanUpload(l.svcCtx, userIdentity); err != nil {
		return nil, err
	}

	req.Name = helper.NormalizeInput(req.Name)
	req.RepositoryIdentity = helper.NormalizeInput(req.RepositoryIdentity)

	if req.RepositoryIdentity == "" {
		return nil, errors.New("资源标识不能为空")
	}
	if req.Name != "" && !helper.IsValidFileOrFolderName(req.Name) {
		return nil, errors.New("名称不合法")
	}

	availableName, err := l.buildAvailableName(userIdentity, req.ParentId, req.Name)
	if err != nil {
		return nil, err
	}

	ur := models.UserRepository{
		Identity:           helper.GenerateUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                req.Ext,
		Name:               availableName,
	}
	_, err = l.svcCtx.Engine.Insert(ur)
	if err != nil {
		return nil, err
	}
	helper.AddAuditLog(l.svcCtx.Engine, userIdentity, "", 1, "FILE_SAVE_REPOSITORY", "user_repository", ur.Identity, "用户保存文件到个人网盘")
	return
}

func (l *UserRepositoryLogic) buildAvailableName(userIdentity string, parentId int64, rawName string) (string, error) {
	return buildAvailableName(l.svcCtx.Engine, userIdentity, parentId, rawName)
}
