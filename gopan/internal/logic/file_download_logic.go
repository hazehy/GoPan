// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopan/gopan/define"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadLogic {
	return &FileDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDownloadLogic) FileDownload(req *types.FileDownloadRequest, userIdentity string) (resp *types.FileDownloadResponse, err error) {
	if err := ensureUserCanDownload(l.svcCtx, userIdentity); err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.RepositoryIdentity) == "" {
		return nil, fmt.Errorf("资源标识不能为空")
	}

	rp := new(models.RepositoryPool)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.RepositoryIdentity).Get(rp)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("资源不存在")
	}

	bucketURL, err := url.Parse(define.COSBucketURL)
	if err != nil {
		return nil, err
	}

	objectKey := strings.TrimPrefix(strings.TrimSpace(pathFromURL(rp.Path)), "/")
	if objectKey == "" {
		return nil, fmt.Errorf("资源路径无效")
	}

	filename := strings.TrimSpace(req.Filename)
	if filename == "" {
		filename = rp.Name + rp.Ext
	}

	c := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{})
	opt := &cos.PresignedURLOptions{Query: &url.Values{}}
	opt.Query.Set("response-content-disposition", contentDisposition(filename))
	if detected := mimeFromExt(rp.Ext); detected != "" {
		opt.Query.Set("response-content-type", detected)
	}

	presignedURL, err := c.Object.GetPresignedURL(
		l.ctx,
		http.MethodGet,
		objectKey,
		define.TencentSecretID,
		define.TencentSecretKey,
		time.Minute*10,
		opt,
		true,
	)
	if err != nil {
		return nil, err
	}

	resp = &types.FileDownloadResponse{Url: presignedURL.String()}
	return resp, nil
}

func pathFromURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	return parsed.Path
}

func contentDisposition(filename string) string {
	escaped := url.QueryEscape(filename)
	quoted := strings.ReplaceAll(filename, `"`, "")
	return fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, quoted, escaped)
}

func mimeFromExt(ext string) string {
	switch strings.ToLower(strings.TrimSpace(ext)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	case ".txt", ".log":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/vnd.rar"
	case ".7z":
		return "application/x-7z-compressed"
	default:
		return ""
	}
}
