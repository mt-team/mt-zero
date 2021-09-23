package handler

import (
	"net/http"

    "ruquan/src/util/response"
	{{.ImportPackages}}
)

func {{.HandlerName}}(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(w, response.ErrInvalidArgs)
			return
		}{{end}}

		l := logic.New{{.LogicType}}(r.Context(), ctx)
		span := c.Value(tracespec.TracingKey).(tracespec.Trace)
        w.Header().Set("X-Trace-ID", span.TraceId())
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}req{{end}})
		if err != nil {
            response.Response(w, err)
            return
        }

        {{if .HasResp}}response.Response(w, resp){{else}}response.Response(w, nil){{end}}
	}
}
