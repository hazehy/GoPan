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

	availableName, err := buildAvailableName(l.svcCtx.Engine, userIdentity, req.ParentId, req.Name)
	if err != nil {
		return nil, err
	}

	data := &models.UserRepository{
		Identity:     helper.GenerateUUID(),
		UserIdentity: userIdentity,
		ParentId:     req.ParentId,
		Name:         availableName,
	}
	_, err = l.svcCtx.Engine.Insert(data)
	if err != nil {
		return nil, err
	}
	resp = &types.FolderCreateResponse{Identity: data.Identity}
	return
}
