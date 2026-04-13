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

type FileRenameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileRenameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileRenameLogic {
	return &FileRenameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileRenameLogic) FileRename(req *types.FileRenameRequest, userIdentity string) (resp *types.FileRenameResponse, err error) {
	req.Name = helper.NormalizeInput(req.Name)
	if !helper.IsValidFileOrFolderName(req.Name) {
		return nil, errors.New("名称不合法")
	}

	cnt, err := l.svcCtx.Engine.Where("name = ? AND parent_id = (SELECT parent_id FROM user_repository WHERE user_repository.identity = ?)", req.Name, req.Identity).Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("文件名已存在")
	}
	data := &models.UserRepository{
		Name: req.Name,
	}
	_, err = l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.Identity, userIdentity).Update(data)
	if err != nil {
		return nil, err
	}
	return
}
