package svc

import (
	"ruquan/src/app/internal/config"
	"ruquan/src/app/internal/model"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	VersionModel model.VersionModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:       c,
		VersionModel: model.NewVersionModel(conn, c.CacheRedis),
	}
}
