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

type AdminUserStatusUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserStatusUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserStatusUpdateLogic {
	return &AdminUserStatusUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserStatusUpdateLogic) AdminUserStatusUpdate(req *types.AdminUserStatusUpdateRequest) error {
	req.Identity = helper.NormalizeInput(req.Identity)
	if req.Identity == "" {
		return errors.New("用户标识不能为空")
	}

	if req.Status == nil && req.UploadPermission == nil && req.DownloadPermission == nil && req.SharePermission == nil {
		return errors.New("至少需要更新一个字段")
	}

	if req.Status != nil && *req.Status != 1 && *req.Status != 2 {
		return errors.New("用户状态仅支持 1(正常) 或 2(禁用)")
	}
	if req.UploadPermission != nil && *req.UploadPermission != 1 && *req.UploadPermission != 2 {
		return errors.New("上传权限仅支持 1(允许) 或 2(禁止)")
	}
	if req.DownloadPermission != nil && *req.DownloadPermission != 1 && *req.DownloadPermission != 2 {
		return errors.New("下载权限仅支持 1(允许) 或 2(禁止)")
	}
	if req.SharePermission != nil && *req.SharePermission != 1 && *req.SharePermission != 2 {
		return errors.New("分享权限仅支持 1(允许) 或 2(禁止)")
	}

	target := new(models.User)
	has, err := l.svcCtx.Engine.Where("identity = ? AND deleted_at IS NULL", req.Identity).Get(target)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("用户不存在")
	}
	if target.Role == 2 {
		return errors.New("不能修改管理员状态")
	}

	updated := &models.User{}
	cols := make([]string, 0, 4)
	auditParts := make([]string, 0, 4)
	hasActualChange := false

	if req.Status != nil {
		if target.Status != *req.Status {
			hasActualChange = true
		}
		updated.Status = *req.Status
		cols = append(cols, "status")
		auditParts = append(auditParts, fmt.Sprintf("status=%d", *req.Status))
	}
	if req.UploadPermission != nil {
		if target.UploadPermission != *req.UploadPermission {
			hasActualChange = true
		}
		updated.UploadPermission = *req.UploadPermission
		cols = append(cols, "upload_permission")
		auditParts = append(auditParts, fmt.Sprintf("upload_permission=%d", *req.UploadPermission))
	}
	if req.DownloadPermission != nil {
		if target.DownloadPermission != *req.DownloadPermission {
			hasActualChange = true
		}
		updated.DownloadPermission = *req.DownloadPermission
		cols = append(cols, "download_permission")
		auditParts = append(auditParts, fmt.Sprintf("download_permission=%d", *req.DownloadPermission))
	}
	if req.SharePermission != nil {
		if target.SharePermission != *req.SharePermission {
			hasActualChange = true
		}
		updated.SharePermission = *req.SharePermission
		cols = append(cols, "share_permission")
		auditParts = append(auditParts, fmt.Sprintf("share_permission=%d", *req.SharePermission))
	}

	if len(cols) == 0 {
		return errors.New("没有有效的更新字段")
	}
	if !hasActualChange {
		return nil
	}

	affected, err := l.svcCtx.Engine.Where("identity = ? AND deleted_at IS NULL", req.Identity).Cols(cols...).Update(updated)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("未更新到任何用户记录")
	}

	helper.AddAuditLog(l.svcCtx.Engine, "SYSTEM", "admin", 2, "USER_STATUS_UPDATE", "user", target.Identity, "管理员更新用户状态与权限: "+strings.Join(auditParts, ";"))
	return nil
}
