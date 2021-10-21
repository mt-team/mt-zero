// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"mtzero/src/gateway/internal/svc"

	"github.com/tal-tech/go-zero/rest"
)

func RegisterHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {
	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Trace},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/app/v1/info",
					Handler: AppInfoHandler(serverCtx),
				},
			}...,
		),
	)
}
