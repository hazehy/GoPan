// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileChunkUploadCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileChunkUploadCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileChunkUploadCompleteLogic {
	return &FileChunkUploadCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileChunkUploadCompleteLogic) FileChunkUploadComplete(req *types.FileChunkUploadCompleteRequest, userIdentity string) (resp *types.FileChunkUploadCompleteResponse, err error) {
	if err := ensureUserCanUpload(l.svcCtx, userIdentity); err != nil {
		return nil, err
	}

	parts := make([]cos.Object, 0)
	for _, part := range req.CosObjects {
		parts = append(parts, cos.Object{
			ETag:       part.Etag,
			PartNumber: part.PartNumber,
		})
	}
	err = helper.CosChunkComplete(req.Key, req.UploadId, parts)
	if err != nil {
		return nil, err
	}
	return
}
