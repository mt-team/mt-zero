package app

import (
	"context"

	"mtzero/src/app/app"
	"mtzero/src/gateway/internal/svc"
	"mtzero/src/gateway/internal/types"
	bizResponse "mtzero/src/util/response"

	"github.com/tal-tech/go-zero/core/logx"
)

type AppInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) AppInfoLogic {
	return AppInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppInfoLogic) AppInfo(req types.AppReq) (*types.AppResp, error) {
	resp, err := l.svcCtx.AppRpc.GetAppInfo(l.ctx, &app.GetAppInfoReq{
		Platform: req.Plt,
		Version:  req.Ver,
	})
	if err != nil {
		l.Logger.Errorf("AppInfo err %v", err)
		return nil, bizResponse.ErrAppCtl
	}

	return &types.AppResp{
		Version: resp.Version,
		Desc:    resp.Desc,
		Url:     resp.Url,
		Type:    int8(resp.Type),
		Env:     resp.Env,
	}, nil
}
