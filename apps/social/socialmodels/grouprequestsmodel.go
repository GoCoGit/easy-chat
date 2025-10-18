package socialmodels

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupRequestsModel = (*customGroupRequestsModel)(nil)

type (
	// GroupRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupRequestsModel.
	GroupRequestsModel interface {
		groupRequestsModel
		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
		UpdateTrans(ctx context.Context, session sqlx.Session, data *GroupRequests) error
		ListNoHandler(ctx context.Context, groupId string) ([]*GroupRequests, error)
		FindByGroupIdAndReqId(ctx context.Context, groupId, reqId string) (*GroupRequests, error)
	}

	customGroupRequestsModel struct {
		*defaultGroupRequestsModel
	}
)

// NewGroupRequestsModel returns a model for the database table.
func NewGroupRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupRequestsModel {
	return &customGroupRequestsModel{
		defaultGroupRequestsModel: newGroupRequestsModel(conn, c, opts...),
	}
}

func (m *defaultGroupRequestsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultGroupRequestsModel) UpdateTrans(ctx context.Context, session sqlx.Session, data *GroupRequests) error {
	groupRequestsIdKey := fmt.Sprintf("%s%v", cacheGroupRequestsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupRequestsRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.ReqId, data.GroupId, data.ReqMsg, data.ReqTime, data.JoinSource,
			data.InviterUserId, data.HandleUserId, data.HandleTime, data.HandleResult, data.Id)
	}, groupRequestsIdKey)
	return err
}

func (m *defaultGroupRequestsModel) ListNoHandler(ctx context.Context, groupId string) ([]*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `handle_result` = 1 ", groupRequestsRows, m.table)

	var resp []*GroupRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultGroupRequestsModel) FindByGroupIdAndReqId(ctx context.Context, groupId, reqId string) (*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `req_id` = ? and `group_id` = ?", groupRequestsRows, m.table)

	var resp GroupRequests
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, reqId, groupId)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}
