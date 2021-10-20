package handler

import (
	"context"
	"net/http"

	"mtzero/src/gateway/internal/logic"
	"mtzero/src/gateway/internal/svc"
	"mtzero/src/gateway/internal/types"
	"mtzero/src/gateway/util"
	bizResponse "mtzero/src/util/response"

	"github.com/tal-tech/go-zero/core/trace/tracespec"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func AppInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AppReq
		if err := httpx.Parse(r, &req); err != nil {
			bizResponse.Response(w, err)
			return
		}

		userUuid := util.GetHeaderUserUuid(r)
		c := context.WithValue(r.Context(), "userUuid", userUuid)
		c = util.CpoyHeaderToCtx(c, r)
		span := c.Value(tracespec.TracingKey).(tracespec.Trace)
		w.Header().Set("X-Trace-ID", span.TraceId())

		l := logic.NewAppInfoLogic(r.Context(), ctx)
		resp, err := l.AppInfo(req)
		if err != nil {
			bizResponse.Response(w, err)
			return
		}

		bizResponse.Response(w, resp)
	}
}
