package handler

import (
	"net/http"

	"ruquan/src/gateway/internal/logic"
	"ruquan/src/gateway/internal/svc"
	"ruquan/src/gateway/internal/types"
	"ruquan/src/util/response"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func AppInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AppReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(w, response.ErrInvalidArgs)
			return
		}

		l := logic.NewAppInfoLogic(r.Context(), ctx)
		resp, err := l.AppInfo(req)
		if err != nil {
			response.Response(w, err)
			return
		}

		response.Response(w, resp)
	}
}
