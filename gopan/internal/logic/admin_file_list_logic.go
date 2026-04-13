package logic

import (
	"context"
	"fmt"
	"strings"

	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminFileListLogic {
	return &AdminFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminFileListLogic) AdminFileList(req *types.AdminFileListRequest) (resp *types.AdminFileListResponse, err error) {
	_, size, offset := normalizePageAndSize(req.Page, req.Size)

	type adminFileRow struct {
		Id                 int64  `xorm:"id"`
		Identity           string `xorm:"identity"`
		ParentId           int64  `xorm:"parent_id"`
		UserIdentity       string `xorm:"user_identity"`
		UserName           string `xorm:"user_name"`
		RepositoryIdentity string `xorm:"repository_identity"`
		Name               string `xorm:"name"`
		Ext                string `xorm:"ext"`
		Path               string `xorm:"path"`
		Size               int64  `xorm:"size"`
		UpdatedAt          string `xorm:"updated_at"`
	}

	rows := make([]*adminFileRow, 0)
	querySession := l.svcCtx.Engine.Table("user_repository").
		Select("user_repository.id, user_repository.identity, user_repository.parent_id, user_repository.user_identity, user_basic.name AS user_name, user_repository.repository_identity, user_repository.name, user_repository.ext, user_repository.updated_at, repository_pool.path, repository_pool.size").
		Join("LEFT", "repository_pool", "user_repository.repository_identity = repository_pool.identity").
		Join("LEFT", "user_basic", "user_repository.user_identity = user_basic.identity").
		Where("user_repository.deleted_at IS NULL")
	countSession := l.svcCtx.Engine.Table("user_repository").
		Join("LEFT", "user_basic", "user_repository.user_identity = user_basic.identity").
		Where("user_repository.deleted_at IS NULL")

	if userName := strings.TrimSpace(req.UserName); userName != "" {
		likeUserName := "%" + userName + "%"
		querySession = querySession.And("user_basic.name LIKE ?", likeUserName)
		countSession = countSession.And("user_basic.name LIKE ?", likeUserName)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		querySession = querySession.And("(user_repository.name LIKE ? OR user_repository.identity LIKE ? OR user_repository.user_identity LIKE ? OR user_basic.name LIKE ?)", like, like, like, like)
		countSession = countSession.And("(user_repository.name LIKE ? OR user_repository.identity LIKE ? OR user_repository.user_identity LIKE ? OR user_basic.name LIKE ?)", like, like, like, like)
	}

	err = querySession.Desc("user_repository.id").Limit(size, offset).Find(&rows)
	if err != nil {
		return nil, err
	}

	pathCache := map[string]string{}
	list := make([]*types.AdminFileItem, 0, len(rows))
	for _, row := range rows {
		dirPath, pathErr := l.buildDirectoryPath(row.UserIdentity, row.ParentId, pathCache)
		if pathErr != nil {
			return nil, pathErr
		}
		list = append(list, &types.AdminFileItem{
			Identity:           row.Identity,
			ParentId:           row.ParentId,
			UserIdentity:       row.UserIdentity,
			UserName:           row.UserName,
			RepositoryIdentity: row.RepositoryIdentity,
			Name:               row.Name,
			Ext:                row.Ext,
			Path:               dirPath,
			Size:               row.Size,
			UpdatedAt:          row.UpdatedAt,
		})
	}

	count, err := countSession.Count(new(models.UserRepository))
	if err != nil {
		return nil, err
	}

	resp = &types.AdminFileListResponse{List: list, Count: count}
	return resp, nil
}

func (l *AdminFileListLogic) buildDirectoryPath(userIdentity string, parentID int64, cache map[string]string) (string, error) {
	if parentID <= 0 {
		return "根目录", nil
	}

	type folderNode struct {
		Id       int64
		ParentId int64
		Name     string
	}

	lineage := make([]folderNode, 0, 8)
	currentID := parentID
	prefix := ""

	for currentID > 0 {
		cacheKey := fmt.Sprintf("%s:%d", userIdentity, currentID)
		if cached, ok := cache[cacheKey]; ok {
			prefix = cached
			break
		}

		node := struct {
			Id       int64  `xorm:"id"`
			ParentId int64  `xorm:"parent_id"`
			Name     string `xorm:"name"`
		}{}
		has, err := l.svcCtx.Engine.Table("user_repository").
			Select("id, parent_id, name").
			Where("id = ? AND user_identity = ? AND deleted_at IS NULL", currentID, userIdentity).
			Get(&node)
		if err != nil {
			return "", err
		}
		if !has {
			break
		}

		lineage = append(lineage, folderNode{
			Id:       node.Id,
			ParentId: node.ParentId,
			Name:     strings.TrimSpace(node.Name),
		})
		currentID = node.ParentId
	}

	segments := make([]string, 0, len(lineage)+4)
	if prefix != "" {
		segments = append(segments, strings.Split(prefix, "/")...)
	}

	for i := len(lineage) - 1; i >= 0; i-- {
		node := lineage[i]
		if node.Name != "" {
			segments = append(segments, node.Name)
		}
		cacheKey := fmt.Sprintf("%s:%d", userIdentity, node.Id)
		cache[cacheKey] = strings.Join(segments, "/")
	}

	if len(segments) == 0 {
		return "根目录", nil
	}
	return "根目录/" + strings.Join(segments, "/"), nil
}
