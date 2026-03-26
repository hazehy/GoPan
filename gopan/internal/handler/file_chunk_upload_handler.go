// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"errors"
	"net/http"

	"gopan/gopan/helper"
	"gopan/gopan/internal/logic"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileChunkUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileChunkUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if r.PostForm.Get("key") == "" {
			httpx.Error(w, errors.New("key是空的"))
			return
		}
		if r.PostForm.Get("upload_id") == "" {
			httpx.Error(w, errors.New("upload_id是空的"))
			return
		}
		if r.PostForm.Get("part_number") == "" {
			httpx.Error(w, errors.New("part_number是空的"))
			return
		}

		partETag, err := helper.CosChunkUpload(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewFileChunkUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileChunkUpload(&req)
		resp = new(types.FileChunkUploadResponse)
		resp.Etag = partETag
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
