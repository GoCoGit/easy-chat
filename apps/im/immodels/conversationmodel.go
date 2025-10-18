package immodels

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationModel.
	ConversationModel interface {
		conversationModel
		ListByConversationIds(ctx context.Context, ids []string) ([]*Conversation, error)
		UpdateMsg(ctx context.Context, chatLog *ChatLog) error
		FindOneByConversationId(ctx context.Context, conversationId string) (*Conversations, error)
	}

	customConversationModel struct {
		*defaultConversationModel
	}
)

// NewConversationModel returns a model for the mongo.
func NewConversationModel(url, db, collection string) ConversationModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customConversationModel{
		defaultConversationModel: newDefaultConversationModel(conn),
	}
}

func (m *defaultConversationModel) ListByConversationIds(ctx context.Context, ids []string) ([]*Conversation, error) {
	var data []*Conversation

	err := m.conn.Find(ctx, &data, bson.M{
		"conversationId": bson.M{
			"$in": ids,
		},
	})
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultConversationModel) UpdateMsg(ctx context.Context, chatLog *ChatLog) error {
	_, err := m.conn.UpdateOne(ctx,
		bson.M{"conversationId": chatLog.ConversationId},
		bson.M{
			// 更新会话总消息数
			"$inc": bson.M{"total": 1},
			"$set": bson.M{"msg": chatLog},
		},
	)
	return err
}

func (m *defaultConversationModel) FindOneByConversationId(ctx context.Context, conversationId string) (*Conversations, error) {
	var data Conversations

	err := m.conn.FindOne(ctx, &data, bson.M{"conversationId": conversationId})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
