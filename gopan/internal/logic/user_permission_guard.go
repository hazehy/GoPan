package logic

import (
	"errors"
	"strings"

	"gopan/gopan/internal/svc"
	"gopan/gopan/models"
)

const (
	permissionEnabled  = 1
	permissionDisabled = 2
)

func ensureUserCanUpload(svcCtx *svc.ServiceContext, userIdentity string) error {
	return ensureUserPermission(svcCtx, userIdentity, "upload")
}

func EnsureUserCanUpload(svcCtx *svc.ServiceContext, userIdentity string) error {
	return ensureUserCanUpload(svcCtx, userIdentity)
}

func ensureUserCanDownload(svcCtx *svc.ServiceContext, userIdentity string) error {
	return ensureUserPermission(svcCtx, userIdentity, "download")
}

func ensureUserCanShare(svcCtx *svc.ServiceContext, userIdentity string) error {
	return ensureUserPermission(svcCtx, userIdentity, "share")
}

func ensureUserPermission(svcCtx *svc.ServiceContext, userIdentity string, permission string) error {
	userIdentity = strings.TrimSpace(userIdentity)
	if userIdentity == "" {
		return errors.New("用户身份缺失")
	}

	user := new(models.User)
	has, err := svcCtx.Engine.Where("identity = ? AND deleted_at IS NULL", userIdentity).Get(user)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("用户不存在")
	}
	if user.Status == permissionDisabled {
		return errors.New("当前用户已被禁用")
	}

	switch permission {
	case "upload":
		if user.UploadPermission == permissionDisabled {
			return errors.New("当前账号无上传权限")
		}
	case "download":
		if user.DownloadPermission == permissionDisabled {
			return errors.New("当前账号无下载权限")
		}
	case "share":
		if user.SharePermission == permissionDisabled {
			return errors.New("当前账号无分享权限")
		}
	}

	return nil
}
