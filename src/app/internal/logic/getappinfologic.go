package logic

import (
	"context"

	"mtzero/src/app/app"
	"mtzero/src/app/internal/svc"
	bizResponse "mtzero/src/util/response"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/utils"
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
	mplt := cpltToMplt(in.Platform)
	if mplt < 0 {
		return nil, bizResponse.ErrInvalidArgs.WithMessage("平台信息不正确")
	}

	lastVer, err := l.svcCtx.VersionModel.FindLastOneByPlt(mplt)
	if err != nil {
		l.Logger.Errorf("GetAppInfo VersionModel.FindLastOneByPlt mplt<%v> err %v", mplt, err)
		return nil, err
	}

	lastForceVer, err := l.svcCtx.VersionModel.FindLastForceOneByPlt(mplt)
	if err != nil {
		l.Logger.Errorf("GetAppInfo VersionModel.FindLastForceOneByPlt mplt<%v> err %v", mplt, err)
		return nil, err
	}

	// 更新逻辑判断
	updateType := app.UpdateType_None
	if utils.CompareVersions(in.Version, "<", lastVer.Ver) {
		updateType = app.UpdateType_Soft
	}
	if utils.CompareVersions(in.Version, "<", lastForceVer.Ver) {
		updateType = app.UpdateType_Force
	}

	return &app.GetAppInfoResp{
		Version: lastVer.Ver,
		Desc:    lastVer.Desc,
		Url:     lastVer.Url,
		Type:    updateType,
		Env:     l.svcCtx.Config.Mode,
	}, nil
}

func cpltToMplt(cplt string) int8 {
	// 0:h5,1:iOS,2:android,3:wxma
	switch cplt {
	case "h5":
		return 0
	case "iOS":
		return 1
	case "android":
		return 2
	case "wxma":
		return 3
	}

	return -1
}
