package svc

import (
	"mtzero/src/app/appclient"
	"mtzero/src/gateway/internal/config"

	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	AppRpc appclient.App
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		AppRpc: appclient.NewApp(zrpc.MustNewClient(c.AppRpc)),
	}
}
