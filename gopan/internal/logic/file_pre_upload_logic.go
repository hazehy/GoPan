// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FilePreUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFilePreUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilePreUploadLogic {
	return &FilePreUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FilePreUploadLogic) FilePreUpload(req *types.FilePreUploadRequest, userIdentity string) (resp *types.FilePreUploadResponse, err error) {
	if err := ensureUserCanUpload(l.svcCtx, userIdentity); err != nil {
		return nil, err
	}

	rp := new(models.RepositoryPool)
	has, err := l.svcCtx.Engine.Where("hash = ?", req.Md5).Get(rp)
	if err != nil {
		return nil, err
	}
	resp = new(types.FilePreUploadResponse)
	if has {
		resp.Identity = rp.Identity
	} else {
		key, uploadId, err := helper.CosChunkInit(req.Ext)
		if err != nil {
			return nil, err
		}
		resp.Key = key
		resp.UploadId = uploadId
	}
	return
}
