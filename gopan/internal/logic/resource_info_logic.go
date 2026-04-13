// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceInfoLogic {
	return &ResourceInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceInfoLogic) ResourceInfo(req *types.ResourceInfoRequest) (resp *types.ResourceInfoResponse, err error) {
	if req.Identity == "" {
		return nil, errors.New("分享标识不能为空")
	}
	_, err = l.svcCtx.Engine.Exec("UPDATE share_link SET click_num = click_num + 1 WHERE identity = ?", req.Identity)
	if err != nil {
		return nil, err
	}
	resp = new(types.ResourceInfoResponse)
	_, err = l.svcCtx.Engine.Table("share_link").
		Select("share_link.repository_identity, repository_pool.name, repository_pool.ext, repository_pool.size, repository_pool.path").
		Join("LEFT", "repository_pool", "share_link.repository_identity = repository_pool.identity").
		Where("share_link.identity = ?", req.Identity).
		Get(resp)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.Engine.Table("user_repository").
		Select("name").
		Where("repository_identity = ?", resp.RepositoryIdentity).
		Get(&resp.Name)
	if err != nil {
		return nil, err
	}
	return
}
