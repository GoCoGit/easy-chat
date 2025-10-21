package msgTransfer

import (
	"context"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/bitmap"
	"easy-chat/pkg/constants"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-queue/kq"
)

var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*baseMsgTransfer
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
	}

	return m
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Infof("MsgReadTransfer ", value)
	var data mq.MsgMarkRead

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 业务处理 - 更新对已读未读的记录
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	// map[string]string - 消息ID -> 已读记录

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ReadRecords:    readRecords,
		ContentType:    constants.ContentMarkRead,
	})
}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)

	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}
	// 处理已读
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}

		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)
		err = m.svcCtx.ChatLogModel.UpdateMarkRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
