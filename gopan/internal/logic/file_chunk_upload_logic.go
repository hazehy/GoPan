// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileChunkUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileChunkUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileChunkUploadLogic {
	return &FileChunkUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileChunkUploadLogic) FileChunkUpload(req *types.FileChunkUploadRequest) (resp *types.FileChunkUploadResponse, err error) {
	return
}
