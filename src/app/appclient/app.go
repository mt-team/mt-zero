// Code generated by goctl. DO NOT EDIT!
// Source: app.proto

//go:generate mockgen -destination ./app_mock.go -package appclient -source $GOFILE

package appclient

import (
	"context"

	"ruquan/src/app/app"

	"github.com/tal-tech/go-zero/zrpc"
)

type (
	GetAppInfoReq  = app.GetAppInfoReq
	GetAppInfoResp = app.GetAppInfoResp

	App interface {
		GetAppInfo(ctx context.Context, in *GetAppInfoReq) (*GetAppInfoResp, error)
	}

	defaultApp struct {
		cli zrpc.Client
	}
)

func NewApp(cli zrpc.Client) App {
	return &defaultApp{
		cli: cli,
	}
}

func (m *defaultApp) GetAppInfo(ctx context.Context, in *GetAppInfoReq) (*GetAppInfoResp, error) {
	client := app.NewAppClient(m.cli.Conn())
	return client.GetAppInfo(ctx, in)
}