// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"gopan/gopan/define"
	"gopan/gopan/helper"
	"gopan/gopan/internal/logic"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"
	"gopan/gopan/models"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if req.Hash != "" && req.Name != "" && req.Size > 0 {
			if req.Path == "" && req.Key != "" {
				req.Path = strings.TrimRight(define.COSBucketURL, "/") + "/" + req.Key
			}
			if req.Path == "" {
				httpx.ErrorCtx(r.Context(), w, errors.New("path是空的"))
				return
			}
			rp := new(models.RepositoryPool)
			has, err := svcCtx.Engine.Where("hash = ?", req.Hash).Get(rp)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
			if has {
				role, _ := strconv.Atoi(r.Header.Get("UserRole"))
				ext := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(rp.Ext), "."))
				if ext == "" {
					ext = "none"
				}
				detail := fmt.Sprintf("file_ext=%s;upload_mode=instant;file_name=%s", ext, rp.Name)
				helper.AddAuditLog(svcCtx.Engine, r.Header.Get("UserIdentity"), r.Header.Get("UserName"), role, "FILE_UPLOAD", "repository_pool", rp.Identity, detail)
				httpx.OkJsonCtx(r.Context(), w, &types.FileUploadResponse{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
				return
			}
			l := logic.NewFileUploadLogic(r.Context(), svcCtx)
			resp, err := l.FileUpload(&req)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
			} else {
				role, _ := strconv.Atoi(r.Header.Get("UserRole"))
				ext := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(resp.Ext), "."))
				if ext == "" {
					ext = "none"
				}
				detail := fmt.Sprintf("file_ext=%s;upload_mode=direct;file_name=%s", ext, resp.Name)
				helper.AddAuditLog(svcCtx.Engine, r.Header.Get("UserIdentity"), r.Header.Get("UserName"), role, "FILE_UPLOAD", "repository_pool", resp.Identity, detail)
				httpx.OkJsonCtx(r.Context(), w, resp)
			}
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("缺少file文件或元数据参数不完整"))
			return
		}
		// 计算文件内容哈希
		b := make([]byte, fileHeader.Size)
		_, err = file.Read(b)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		hash := fmt.Sprintf("%x", sha256.Sum256(b))
		rp := new(models.RepositoryPool)
		// 判断文件是否存在
		has, err := svcCtx.Engine.Where("hash = ?", hash).Get(rp)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if has {
			role, _ := strconv.Atoi(r.Header.Get("UserRole"))
			ext := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(rp.Ext), "."))
			if ext == "" {
				ext = "none"
			}
			detail := fmt.Sprintf("file_ext=%s;upload_mode=instant;file_name=%s", ext, rp.Name)
			helper.AddAuditLog(svcCtx.Engine, r.Header.Get("UserIdentity"), r.Header.Get("UserName"), role, "FILE_UPLOAD", "repository_pool", rp.Identity, detail)
			httpx.OkJsonCtx(r.Context(), w, &types.FileUploadResponse{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
			return
		}
		// 上传文件到COS
		COSPath, err := helper.COSUpLoad(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 传递参数到logic层
		req.Name = fileHeader.Filename
		req.Ext = path.Ext(fileHeader.Filename)
		req.Size = fileHeader.Size
		req.Hash = hash
		req.Path = COSPath

		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			role, _ := strconv.Atoi(r.Header.Get("UserRole"))
			ext := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(resp.Ext), "."))
			if ext == "" {
				ext = "none"
			}
			detail := fmt.Sprintf("file_ext=%s;upload_mode=direct;file_name=%s", ext, resp.Name)
			helper.AddAuditLog(svcCtx.Engine, r.Header.Get("UserIdentity"), r.Header.Get("UserName"), role, "FILE_UPLOAD", "repository_pool", resp.Identity, detail)
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
