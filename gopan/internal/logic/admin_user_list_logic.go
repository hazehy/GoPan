package logic

import (
	"context"
	"strings"
	"time"

	"gopan/gopan/define"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserListLogic {
	return &AdminUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserListLogic) AdminUserList(req *types.AdminUserListRequest) (resp *types.AdminUserListResponse, err error) {
	_, size, offset := normalizePageAndSize(req.Page, req.Size)

	type userRow struct {
		Identity           string    `xorm:"identity"`
		Name               string    `xorm:"name"`
		Email              string    `xorm:"email"`
		Status             int       `xorm:"status"`
		Role               int       `xorm:"role"`
		UploadPermission   int       `xorm:"upload_permission"`
		DownloadPermission int       `xorm:"download_permission"`
		SharePermission    int       `xorm:"share_permission"`
		LastLoginAt        time.Time `xorm:"last_login_at"`
		CreatedAt          time.Time `xorm:"created_at"`
	}

	rows := make([]*userRow, 0)
	querySession := l.svcCtx.Engine.Table("user_basic").Where("deleted_at IS NULL")
	countSession := l.svcCtx.Engine.Table("user_basic").Where("deleted_at IS NULL")

	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		querySession = querySession.And("(name LIKE ? OR email LIKE ? OR identity LIKE ?)", like, like, like)
		countSession = countSession.And("(name LIKE ? OR email LIKE ? OR identity LIKE ?)", like, like, like)
	}

	err = querySession.
		Select("identity, name, email, status, role, upload_permission, download_permission, share_permission, last_login_at, created_at").
		Desc("id").
		Limit(size, offset).
		Find(&rows)
	if err != nil {
		return nil, err
	}

	count, err := countSession.Count(new(models.User))
	if err != nil {
		return nil, err
	}

	list := make([]*types.AdminUserItem, 0, len(rows))
	for _, row := range rows {
		item := &types.AdminUserItem{
			Identity:           row.Identity,
			Name:               row.Name,
			Email:              row.Email,
			Status:             row.Status,
			Role:               row.Role,
			UploadPermission:   row.UploadPermission,
			DownloadPermission: row.DownloadPermission,
			SharePermission:    row.SharePermission,
		}
		if !row.LastLoginAt.IsZero() {
			item.LastLoginAt = row.LastLoginAt.Format(define.DateFormat)
		}
		if !row.CreatedAt.IsZero() {
			item.CreatedAt = row.CreatedAt.Format(define.DateFormat)
		}
		list = append(list, item)
	}

	resp = &types.AdminUserListResponse{List: list, Count: count}
	return resp, nil
}
