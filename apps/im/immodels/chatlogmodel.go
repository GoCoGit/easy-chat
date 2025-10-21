package immodels

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ ChatLogModel = (*customChatLogModel)(nil)

var DefaultChatLogLimit int64 = 100

type (
	// ChatLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatLogModel.
	ChatLogModel interface {
		chatLogModel
		ListBySendTime(ctx context.Context, conversationId string, startSendTime, endSendTime, limit int64) ([]*ChatLog, error)
		ListByMsgIds(ctx context.Context, ids []string) ([]*ChatLog, error)
		UpdateMarkRead(ctx context.Context, id bson.ObjectID, readRecords []byte) error
		InsertWithId(ctx context.Context, data *ChatLog) error
	}

	customChatLogModel struct {
		*defaultChatLogModel
	}
)

func (m *customChatLogModel) ListBySendTime(ctx context.Context, conversationId string, startSendTime, endSendTime, limit int64) ([]*ChatLog, error) {
	var data []*ChatLog

	opt := options.Find().SetLimit(DefaultChatLogLimit).SetSort(bson.M{
		"sendTime": -1,
	})
	if limit > 0 {
		opt.SetLimit(limit)
	}

	filter := bson.M{
		"conversationId": conversationId,
	}

	if endSendTime > 0 {
		filter["sendTime"] = bson.M{
			"$gt":  endSendTime,
			"$lte": startSendTime,
		}
	} else {
		filter["sendTime"] = bson.M{
			"$lt": startSendTime,
		}
	}
	err := m.conn.Find(ctx, &data, filter, opt)
	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewChatLogModel returns a model for the mongo.
func NewChatLogModel(url, db, collection string) ChatLogModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customChatLogModel{
		defaultChatLogModel: newDefaultChatLogModel(conn),
	}
}

func (m *defaultChatLogModel) ListByMsgIds(ctx context.Context, msgIds []string) ([]*ChatLog, error) {
	var data []*ChatLog
	ids := make([]bson.ObjectID, 0, len(msgIds))
	for _, id := range msgIds {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			fmt.Printf("Failed to convert id %s to ObjectID: %v\n", id, err)
			continue
		}
		ids = append(ids, oid)
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": ids, // 注意这里应该是 $in 而不是 &in
		},
	}

	err := m.conn.Find(ctx, &data, filter)

	switch err {
	case nil:
		return data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultChatLogModel) UpdateMarkRead(ctx context.Context, id bson.ObjectID, readRecords []byte) error {
	_, err := m.conn.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"readRecords": readRecords,
	}})
	return err
}

func (m *defaultChatLogModel) InsertWithId(ctx context.Context, data *ChatLog) error {
	// if data.ID.IsZero() {
	// 	data.ID = bson.NewObjectID()
	// 	data.CreateAt = time.Now()
	// 	data.UpdateAt = time.Now()
	// }

	_, err := m.conn.InsertOne(ctx, data)
	return err
}
