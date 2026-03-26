// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"gopan/gopan/internal/logic"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FolderCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FolderCreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewFolderCreateLogic(r.Context(), svcCtx)
		resp, err := l.FolderCreate(&req, r.Header.Get("userIdentity"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
