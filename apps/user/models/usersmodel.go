package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		FindByPhone(ctx context.Context, phone string) (*Users, error)
		ListByName(ctx context.Context, name string) ([]*Users, error)
		ListByIds(ctx context.Context, ids []string) ([]*Users, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

func (m *customUsersModel) FindByPhone(ctx context.Context, phone string) (*Users, error) {
	usersIdKey := fmt.Sprintf("%s%v", cacheUsersIdPrefix, phone)
	var resp Users
	err := m.QueryRowCtx(ctx, &resp, usersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", usersRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, phone)
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

// ListByName 根据用户名模糊查询用户列表
func (m *customUsersModel) ListByName(ctx context.Context, name string) ([]*Users, error) {
	var resp []*Users
	query := fmt.Sprintf("select %s from %s where `nickname` like ?", usersRows, m.table)
	likeName := "%" + name + "%"
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, likeName)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return resp, nil
}

func (m *customUsersModel) ListByIds(ctx context.Context, ids []string) ([]*Users, error) {
	var resp []*Users
	query := fmt.Sprintf("select %s from %s where `id` in ('%s')", usersRows, m.table, strings.Join(ids, "','"))
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return resp, nil
}

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c, opts...),
	}
}
