package svc

import (
	"mtzero/src/app/internal/config"
	"mtzero/src/app/internal/model"

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
