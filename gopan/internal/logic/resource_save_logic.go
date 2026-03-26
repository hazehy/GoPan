// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gopan/gopan/helper"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceSaveLogic {
	return &ResourceSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceSaveLogic) ResourceSave(req *types.ResourceSaveRequest, userIdentity string) (resp *types.ResourceSaveResponse, err error) {
	rp := new(models.RepositoryPool)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.RepositoryIdentity).Get(rp)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("资源不存在")
	}

	availableName, err := l.buildAvailableName(userIdentity, req.ParentId, rp.Name)
	if err != nil {
		return nil, err
	}

	ur := &models.UserRepository{
		Identity:           helper.GenerateUUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		RepositoryIdentity: req.RepositoryIdentity,
		Ext:                rp.Ext,
		Name:               availableName,
	}
	_, err = l.svcCtx.Engine.Insert(ur)
	if err != nil {
		return nil, err
	}

	saver := new(models.User)
	_, _ = l.svcCtx.Engine.Where("identity = ?", userIdentity).Get(saver)
	saverName := strings.TrimSpace(saver.Name)
	if saverName == "" {
		saverName = userIdentity
	}

	sharerName := "unknown"
	if shareIdentity := strings.TrimSpace(req.ShareIdentity); shareIdentity != "" {
		shareRow := struct {
			SharerName string `xorm:"sharer_name"`
		}{}
		hasShare, shareErr := l.svcCtx.Engine.Table("share_link").
			Select("user_basic.name AS sharer_name").
			Join("LEFT", "user_basic", "share_link.user_identity = user_basic.identity").
			Where("share_link.identity = ? AND share_link.repository_identity = ?", shareIdentity, req.RepositoryIdentity).
			Get(&shareRow)
		if shareErr == nil && hasShare && strings.TrimSpace(shareRow.SharerName) != "" {
			sharerName = strings.TrimSpace(shareRow.SharerName)
		}
	}

	fileExt := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(rp.Ext), "."))
	if fileExt == "" {
		fileExt = "none"
	}
	detail := fmt.Sprintf("sharer_name=%s;saver_name=%s;file_ext=%s;saved_name=%s", sharerName, saverName, fileExt, availableName)
	saverRole := saver.Role
	if saverRole <= 0 {
		saverRole = 1
	}
	helper.AddAuditLog(l.svcCtx.Engine, userIdentity, saverName, saverRole, "SHARE_RESOURCE_SAVE", "user_repository", ur.Identity, detail)
	resp = new(types.ResourceSaveResponse)
	resp.Identity = ur.Identity
	return
}

func (l *ResourceSaveLogic) buildAvailableName(userIdentity string, parentId int64, rawName string) (string, error) {
	baseName := rawName
	if baseName == "" {
		baseName = "未命名文件"
	}

	count, err := l.svcCtx.Engine.Where("user_identity = ? AND parent_id = ? AND name = ?", userIdentity, parentId, baseName).Count(new(models.UserRepository))
	if err != nil {
		return "", err
	}
	if count == 0 {
		return baseName, nil
	}

	for index := 1; ; index++ {
		candidate := fmt.Sprintf("%s(%d)", baseName, index)
		candidateCount, candidateErr := l.svcCtx.Engine.Where("user_identity = ? AND parent_id = ? AND name = ?", userIdentity, parentId, candidate).Count(new(models.UserRepository))
		if candidateErr != nil {
			return "", candidateErr
		}
		if candidateCount == 0 {
			return candidate, nil
		}
	}
}
