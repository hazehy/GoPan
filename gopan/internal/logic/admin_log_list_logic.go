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

type AdminLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLogListLogic {
	return &AdminLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLogListLogic) AdminLogList(req *types.AdminLogListRequest) (resp *types.AdminLogListResponse, err error) {
	size := req.Size
	if size <= 0 {
		size = define.PageSize
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * size

	querySession := l.svcCtx.Engine.Table("audit_log").Where("deleted_at IS NULL")
	countSession := l.svcCtx.Engine.Table("audit_log").Where("deleted_at IS NULL")

	if action := strings.TrimSpace(req.Action); action != "" {
		querySession = querySession.And("action = ?", action)
		countSession = countSession.And("action = ?", action)
	}
	if actorName := strings.TrimSpace(req.ActorName); actorName != "" {
		likeActor := "%" + actorName + "%"
		querySession = querySession.And("actor_name LIKE ?", likeActor)
		countSession = countSession.And("actor_name LIKE ?", likeActor)
	}
	if fileExt := strings.TrimSpace(req.FileExt); fileExt != "" {
		normalizedExt := strings.ToLower(strings.TrimPrefix(fileExt, "."))
		if normalizedExt != "" {
			detailLike := "%file_ext=" + normalizedExt + "%"
			querySession = querySession.And("detail LIKE ?", detailLike)
			countSession = countSession.And("detail LIKE ?", detailLike)
		}
	}
	if sharer := strings.TrimSpace(req.Sharer); sharer != "" {
		likeSharer := "%" + sharer + "%"
		querySession = querySession.And("(detail LIKE ? OR actor_name LIKE ?)", "%sharer_name=%"+sharer+"%", likeSharer)
		countSession = countSession.And("(detail LIKE ? OR actor_name LIKE ?)", "%sharer_name=%"+sharer+"%", likeSharer)
	}
	if saver := strings.TrimSpace(req.Saver); saver != "" {
		likeSaver := "%" + saver + "%"
		querySession = querySession.And("(detail LIKE ? OR actor_name LIKE ?)", "%saver_name=%"+saver+"%", likeSaver)
		countSession = countSession.And("(detail LIKE ? OR actor_name LIKE ?)", "%saver_name=%"+saver+"%", likeSaver)
	}
	if day := strings.TrimSpace(req.Day); day != "" {
		parsedDay, parseErr := time.Parse("2006-01-02", day)
		if parseErr == nil {
			dayText := parsedDay.Format("2006-01-02")
			querySession = querySession.And("DATE(created_at) = ?", dayText)
			countSession = countSession.And("DATE(created_at) = ?", dayText)
		}
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		querySession = querySession.And("(actor_name LIKE ? OR detail LIKE ? OR target_identity LIKE ? OR action LIKE ?)", like, like, like, like)
		countSession = countSession.And("(actor_name LIKE ? OR detail LIKE ? OR target_identity LIKE ? OR action LIKE ?)", like, like, like, like)
	}

	rows := make([]*models.AuditLog, 0)
	err = querySession.Desc("id").Limit(size, offset).Find(&rows)
	if err != nil {
		return nil, err
	}

	count, err := countSession.Count(new(models.AuditLog))
	if err != nil {
		return nil, err
	}

	list := make([]*types.AdminLogItem, 0, len(rows))
	for _, item := range rows {
		list = append(list, &types.AdminLogItem{
			Identity:       item.Identity,
			ActorIdentity:  item.ActorIdentity,
			ActorName:      item.ActorName,
			ActorRole:      item.ActorRole,
			Action:         item.Action,
			TargetType:     item.TargetType,
			TargetIdentity: item.TargetIdentity,
			Detail:         item.Detail,
			CreatedAt:      item.CreatedAt.Format(define.DateFormat),
		})
	}

	resp = &types.AdminLogListResponse{List: list, Count: count}
	return resp, nil
}
