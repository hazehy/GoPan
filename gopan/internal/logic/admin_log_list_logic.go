package logic

import (
	"context"
	"errors"
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
	_, size, offset := normalizePageAndSize(req.Page, req.Size)

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
			querySession = querySession.And("LOWER(detail) LIKE ?", detailLike)
			countSession = countSession.And("LOWER(detail) LIKE ?", detailLike)
		}
	}
	if sharer := strings.TrimSpace(req.Sharer); sharer != "" {
		normalizedSharer := strings.ToLower(sharer)
		likeSharer := "%" + sharer + "%"
		querySession = querySession.And("(LOWER(detail) LIKE ? OR actor_name LIKE ?)", "%sharer_name=%"+normalizedSharer+"%", likeSharer)
		countSession = countSession.And("(LOWER(detail) LIKE ? OR actor_name LIKE ?)", "%sharer_name=%"+normalizedSharer+"%", likeSharer)
	}
	if saver := strings.TrimSpace(req.Saver); saver != "" {
		normalizedSaver := strings.ToLower(saver)
		likeSaver := "%" + saver + "%"
		querySession = querySession.And("(LOWER(detail) LIKE ? OR actor_name LIKE ?)", "%saver_name=%"+normalizedSaver+"%", likeSaver)
		countSession = countSession.And("(LOWER(detail) LIKE ? OR actor_name LIKE ?)", "%saver_name=%"+normalizedSaver+"%", likeSaver)
	}
	if day := strings.TrimSpace(req.Day); day != "" {
		if _, parseErr := time.Parse("2006-01-02", day); parseErr != nil {
			return nil, errors.New("day 参数格式应为 YYYY-MM-DD")
		}
		querySession = querySession.And("DATE(created_at) = ?", day)
		countSession = countSession.And("DATE(created_at) = ?", day)
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
