package svc

import (
	"mtzero/src/app/appclient"
	"mtzero/src/gateway/internal/config"
	"mtzero/src/gateway/internal/middleware"

	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	AppRpc appclient.App
	Trace  rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		AppRpc: appclient.NewApp(zrpc.MustNewClient(c.AppRpc)),
		Trace:  middleware.NewTraceMiddleware(c).Handle,
	}
}
