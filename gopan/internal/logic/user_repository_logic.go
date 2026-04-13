// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"

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
	baseName := rawName
	if baseName == "" {
		baseName = "未命名文件"
	}

	count, err := l.svcCtx.Engine.Where("user_identity = ? AND parent_id = ? AND name = ?", userIdentity, parentId, baseName).Count(new(models.UserRepository))
	if err != nil {
		return "", err
	}
	if count == 0 {
		return baseName, nil
	}

	for index := 1; ; index++ {
		candidate := fmt.Sprintf("%s(%d)", baseName, index)
		candidateCount, candidateErr := l.svcCtx.Engine.Where("user_identity = ? AND parent_id = ? AND name = ?", userIdentity, parentId, candidate).Count(new(models.UserRepository))
		if candidateErr != nil {
			return "", candidateErr
		}
		if candidateCount == 0 {
			return candidate, nil
		}
	}
}
