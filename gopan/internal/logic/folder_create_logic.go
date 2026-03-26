// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FolderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FolderCreateLogic {
	return &FolderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FolderCreateLogic) FolderCreate(req *types.FolderCreateRequest, userIdentity string) (resp *types.FolderCreateResponse, err error) {
	req.Name = helper.NormalizeInput(req.Name)
	if !helper.IsValidFileOrFolderName(req.Name) {
		return nil, errors.New("文件夹名称不合法")
	}

	existing := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Table("user_repository").Unscoped().
		Where("user_identity = ? AND parent_id = ? AND name = ?", userIdentity, req.ParentId, req.Name).
		Desc("id").Get(existing)
	if err != nil {
		return nil, err
	}

	if has {
		if existing.DeletedAt.IsZero() {
			return nil, errors.New("文件夹名已存在")
		}

		recycledName := fmt.Sprintf("%s__deleted_%d_%d", req.Name, existing.Id, time.Now().Unix())
		_, err = l.svcCtx.Engine.Exec(
			"UPDATE user_repository SET name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NOT NULL",
			recycledName, time.Now(), existing.Id,
		)
		if err != nil {
			return nil, err
		}
	}

	data := &models.UserRepository{
		Identity:     helper.GenerateUUID(),
		UserIdentity: userIdentity,
		ParentId:     req.ParentId,
		Name:         req.Name,
	}
	_, err = l.svcCtx.Engine.Insert(data)
	if err != nil {
		return nil, err
	}
	resp = &types.FolderCreateResponse{Identity: data.Identity}
	return
}
