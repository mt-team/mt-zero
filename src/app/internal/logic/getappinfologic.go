package logic

import (
	"context"

	"ruquan/src/app/app"
	"ruquan/src/app/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetAppInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppInfoLogic {
	return &GetAppInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppInfoLogic) GetAppInfo(in *app.GetAppInfoReq) (*app.GetAppInfoResp, error) {
	// todo: add your logic here and delete this line

	return &app.GetAppInfoResp{}, nil
}
