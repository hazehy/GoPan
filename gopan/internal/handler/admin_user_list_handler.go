package handler

import (
	"net/http"

	"gopan/gopan/internal/logic"
	"gopan/gopan/internal/svc"
	"gopan/gopan/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminUserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUserListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAdminUserListLogic(r.Context(), svcCtx)
		resp, err := l.AdminUserList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
