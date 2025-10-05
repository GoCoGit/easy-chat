package socialmodels

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendRequestsModel = (*customFriendRequestsModel)(nil)

type (
	// FriendRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendRequestsModel.
	FriendRequestsModel interface {
		friendRequestsModel
		FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error)
		Trans(ctx context.Context, fn func(ctx context.Context,
			session sqlx.Session) error) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *FriendRequests) (sql.Result, error)
		UpdateTrans(ctx context.Context, session sqlx.Session, data *FriendRequests) error
	}

	customFriendRequestsModel struct {
		*defaultFriendRequestsModel
	}
)

// NewFriendRequestsModel returns a model for the database table.
func NewFriendRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendRequestsModel {
	return &customFriendRequestsModel{
		defaultFriendRequestsModel: newFriendRequestsModel(conn, c, opts...),
	}
}

func (m *defaultFriendRequestsModel) FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `user_id` = ?", friendRequestsRows, m.table)

	var resp FriendRequests
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, rid, uid)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFriendRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context,
	session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultFriendRequestsModel) InsertTrans(ctx context.Context, session sqlx.Session, data *FriendRequests) (sql.Result, error) {
	friendRequestsIdKey := fmt.Sprintf("%s%v", cacheFriendRequestsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, friendRequestsRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime, data.HandleResult, data.HandleMsg, data.HandledAt)
	}, friendRequestsIdKey)
	return ret, err
}

func (m *defaultFriendRequestsModel) UpdateTrans(ctx context.Context, session sqlx.Session, data *FriendRequests) error {
	friendRequestsIdKey := fmt.Sprintf("%s%v", cacheFriendRequestsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, friendRequestsRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime, data.HandleResult, data.HandleMsg, data.HandledAt, data.Id)
	}, friendRequestsIdKey)
	return err
}
