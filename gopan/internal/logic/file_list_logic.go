// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"gopan/gopan/define"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileListLogic) FileList(req *types.FileListRequest, userIdentity string) (resp *types.FileListResponse, err error) {
	fl := make([]*types.UserFile, 0)
	resp = new(types.FileListResponse)
	// 计算分页参数
	size := req.Size
	if size == 0 {
		size = define.PageSize
	}
	page := req.Page
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * size
	// 查询文件列表
	err = l.svcCtx.Engine.Table("user_repository").Where("user_identity = ? AND parent_id = ?", userIdentity, req.Id).
		Select("user_repository.id, user_repository.identity, user_repository.repository_identity, user_repository.ext, "+
			"user_repository.name, user_repository.updated_at, repository_pool.path, repository_pool.size").
		Join("LEFT", "repository_pool", "user_repository.repository_identity = repository_pool.identity").
		Where("user_repository.deleted_at = ? OR user_repository.deleted_at IS NULL", time.Time{}.Format(define.DateFormat)).
		Limit(size, offset).Find(&fl)
	if err != nil {
		return nil, err
	}

	for _, item := range fl {
		if item.RepositoryIdentity == "" {
			folderSize, calcErr := l.calcFolderSize(userIdentity, item.Id)
			if calcErr != nil {
				return nil, calcErr
			}
			item.Size = folderSize
		}
	}
	// 查询文件总数
	cnt, err := l.svcCtx.Engine.Table("user_repository").Where("user_identity = ? AND parent_id = ?", userIdentity, req.Id).Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}
	resp.List = fl
	resp.Count = cnt
	return
}

func (l *FileListLogic) calcFolderSize(userIdentity string, folderId int64) (int64, error) {
	type childRecord struct {
		Id                 int64  `xorm:"id"`
		RepositoryIdentity string `xorm:"repository_identity"`
		Size               int64  `xorm:"size"`
	}

	children := make([]*childRecord, 0)
	err := l.svcCtx.Engine.Table("user_repository").
		Select("user_repository.id, user_repository.repository_identity, repository_pool.size").
		Join("LEFT", "repository_pool", "user_repository.repository_identity = repository_pool.identity").
		Where("user_repository.user_identity = ? AND user_repository.parent_id = ?", userIdentity, folderId).
		Where("user_repository.deleted_at = ? OR user_repository.deleted_at IS NULL", time.Time{}.Format(define.DateFormat)).
		Find(&children)
	if err != nil {
		return 0, err
	}

	var totalSize int64
	for _, child := range children {
		if child.RepositoryIdentity == "" {
			subSize, subErr := l.calcFolderSize(userIdentity, child.Id)
			if subErr != nil {
				return 0, subErr
			}
			totalSize += subSize
			continue
		}
		totalSize += child.Size
	}

	return totalSize, nil
}
