package logic

import (
	"context"
	"strings"

	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminOverviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminOverviewLogic {
	return &AdminOverviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminOverviewLogic) AdminOverview() (resp *types.AdminOverviewResponse, err error) {
	resp = new(types.AdminOverviewResponse)

	if resp.TotalUsers, err = l.svcCtx.Engine.Table("user_basic").Where("deleted_at IS NULL").Count(new(models.User)); err != nil {
		return nil, err
	}
	if resp.ActiveUsers, err = l.svcCtx.Engine.Table("user_basic").Where("status = ? AND deleted_at IS NULL", 1).Count(new(models.User)); err != nil {
		return nil, err
	}
	if resp.DisabledUsers, err = l.svcCtx.Engine.Table("user_basic").Where("status = ? AND deleted_at IS NULL", 2).Count(new(models.User)); err != nil {
		return nil, err
	}
	if resp.TotalFiles, err = l.svcCtx.Engine.Table("user_repository").Where("repository_identity <> '' AND deleted_at IS NULL").Count(new(models.UserRepository)); err != nil {
		return nil, err
	}
	if resp.TotalFolders, err = l.svcCtx.Engine.Table("user_repository").Where("(repository_identity = '' OR repository_identity IS NULL) AND deleted_at IS NULL").Count(new(models.UserRepository)); err != nil {
		return nil, err
	}
	if resp.TodayUploads, err = l.svcCtx.Engine.Table("user_repository").Where("repository_identity <> '' AND DATE(created_at) = CURDATE() AND deleted_at IS NULL").Count(new(models.UserRepository)); err != nil {
		return nil, err
	}
	if resp.TodayRegisters, err = l.svcCtx.Engine.Table("user_basic").Where("DATE(created_at) = CURDATE() AND deleted_at IS NULL").Count(new(models.User)); err != nil {
		return nil, err
	}

	type totalSizeRow struct {
		Total int64 `xorm:"total"`
	}
	sizeRow := new(totalSizeRow)
	has, queryErr := l.svcCtx.Engine.Table("user_repository").
		Join("LEFT", "repository_pool", "user_repository.repository_identity = repository_pool.identity").
		Select("COALESCE(SUM(repository_pool.size), 0) AS total").
		Where("user_repository.repository_identity <> '' AND user_repository.deleted_at IS NULL").
		Get(sizeRow)
	if queryErr != nil {
		return nil, queryErr
	}
	if has {
		resp.TotalFileSize = sizeRow.Total
	}

	type extStatRow struct {
		Ext   string `xorm:"ext"`
		Count int64  `xorm:"count"`
	}

	extRows := make([]*extStatRow, 0)
	err = l.svcCtx.Engine.Table("user_repository").
		Select("LOWER(TRIM(ext)) AS ext, COUNT(1) AS count").
		Where("repository_identity <> '' AND deleted_at IS NULL").
		GroupBy("LOWER(TRIM(ext))").
		Desc("count").
		Limit(8, 0).
		Find(&extRows)
	if err != nil {
		return nil, err
	}

	resp.ExtStats = make([]*types.AdminExtStat, 0, len(extRows))
	for _, item := range extRows {
		ext := strings.TrimSpace(item.Ext)
		if ext == "" {
			ext = "unknown"
		}
		resp.ExtStats = append(resp.ExtStats, &types.AdminExtStat{
			Ext:   ext,
			Count: item.Count,
		})
	}

	return resp, nil
}
