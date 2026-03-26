// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDeleteLogic {
	return &FileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDeleteLogic) FileDelete(req *types.FileDeleteRequest, userIdentity string) (resp *types.FileDeleteResponse, err error) {
	target := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.Identity, userIdentity).Get(target)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("文件或文件夹不存在")
	}

	ids, err := l.collectDeleteIds(userIdentity, int64(target.Id))
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.Engine.In("id", ids).Delete(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	return
}

func (l *FileDeleteLogic) collectDeleteIds(userIdentity string, rootId int64) ([]int64, error) {
	ids := []int64{rootId}

	children := make([]*models.UserRepository, 0)
	err := l.svcCtx.Engine.Where("user_identity = ? AND parent_id = ?", userIdentity, rootId).Find(&children)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		childIds, childErr := l.collectDeleteIds(userIdentity, int64(child.Id))
		if childErr != nil {
			return nil, childErr
		}
		ids = append(ids, childIds...)
	}

	return ids, nil
}
