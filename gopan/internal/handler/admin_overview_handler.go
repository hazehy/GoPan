package handler

import (
	"net/http"

	"gopan/gopan/internal/logic"
	"gopan/gopan/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminOverviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminOverviewLogic(r.Context(), svcCtx)
		resp, err := l.AdminOverview()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
