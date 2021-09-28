package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	versionFieldNames          = builderx.RawFieldNames(&Version{})
	versionRows                = strings.Join(versionFieldNames, ",")
	versionRowsExpectAutoSet   = strings.Join(stringx.Remove(versionFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	versionRowsWithPlaceHolder = strings.Join(stringx.Remove(versionFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheVersionIdPrefix       = "cache:version:id:"
	cacheVersionPltPrefix      = "cache:version:plt:"
	cacheVersionPltForcePrefix = "cache:version:force:plt:"
)

type (
	VersionModel interface {
		Insert(data Version) (sql.Result, error)
		FindOne(id int64) (*Version, error)
		FindLastForceOneByPlt(plt int8) (*Version, error)
		FindLastOneByPlt(plt int8) (*Version, error)
		Update(data Version) error
		Delete(id int64) error
	}

	defaultVersionModel struct {
		sqlc.CachedConn
		table string
	}

	Version struct {
		Id         int64     `db:"id"`
		Platform   int64     `db:"platform"`    // 0:h5,1:iOS,2:android,3:wxma
		Ver        string    `db:"ver"`         // 版本号，如 1.9.1
		Force      int64     `db:"force"`       // 这个版本是否强更，1为需要强更
		IsDeleted  int64     `db:"is_deleted"`  // 0 未删除，1 已删除
		Desc       string    `db:"desc"`        // 版本更新内容
		Url        string    `db:"url"`         // cdn更新地址
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 更新时间
	}
)

func NewVersionModel(conn sqlx.SqlConn, c cache.CacheConf) VersionModel {
	return &defaultVersionModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`version`",
	}
}

func (m *defaultVersionModel) Insert(data Version) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, versionRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.Platform, data.Ver, data.Force, data.IsDeleted, data.Desc, data.Url)

	return ret, err
}

func (m *defaultVersionModel) FindOne(id int64) (*Version, error) {
	versionIdKey := fmt.Sprintf("%s%v", cacheVersionIdPrefix, id)
	var resp Version
	err := m.QueryRow(&resp, versionIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", versionRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultVersionModel) Update(data Version) error {
	versionIdKey := fmt.Sprintf("%s%v", cacheVersionIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, versionRowsWithPlaceHolder)
		return conn.Exec(query, data.Platform, data.Ver, data.Force, data.IsDeleted, data.Desc, data.Url, data.Id)
	}, versionIdKey)
	return err
}

func (m *defaultVersionModel) Delete(id int64) error {

	versionIdKey := fmt.Sprintf("%s%v", cacheVersionIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, versionIdKey)
	return err
}

func (m *defaultVersionModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheVersionIdPrefix, primary)
}

func (m *defaultVersionModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", versionRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultVersionModel) FindLastOneByPlt(plt int8) (*Version, error) {
	versionIdKey := fmt.Sprintf("%s%v", cacheVersionPltPrefix, plt)
	var resp Version
	err := m.QueryRow(&resp, versionIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `platform` = ? order by id desc limit 1", versionRows, m.table)
		return conn.QueryRow(v, query, plt)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultVersionModel) FindLastForceOneByPlt(plt int8) (*Version, error) {
	versionIdKey := fmt.Sprintf("%s%v", cacheVersionPltForcePrefix, plt)
	var resp Version
	err := m.QueryRow(&resp, versionIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `platform` = ? and `force`=1 order by id desc limit 1", versionRows, m.table)
		return conn.QueryRow(v, query, plt)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
