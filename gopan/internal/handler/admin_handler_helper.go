package handler

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func writeLogicJSON(ctx context.Context, w http.ResponseWriter, resp any, err error) {
	if err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}
	httpx.OkJsonCtx(ctx, w, resp)
}
