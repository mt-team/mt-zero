package logic

import (
	"context"

	"ruquan/src/gateway/internal/svc"
	"ruquan/src/gateway/internal/types"

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
	// todo: add your logic here and delete this line

	return &types.AppResp{}, nil
}
