package logic

import (
	"context"
	"errors"

	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminFileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminFileDeleteLogic {
	return &AdminFileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminFileDeleteLogic) AdminFileDelete(req *types.AdminFileDeleteRequest) error {
	target := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.Identity).Get(target)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("文件或文件夹不存在")
	}

	ids, err := l.collectDeleteIds(target.UserIdentity, int64(target.Id))
	if err != nil {
		return err
	}

	_, err = l.svcCtx.Engine.In("id", ids).Delete(new(models.UserRepository))
	return err
}

func (l *AdminFileDeleteLogic) collectDeleteIds(userIdentity string, rootId int64) ([]int64, error) {
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
