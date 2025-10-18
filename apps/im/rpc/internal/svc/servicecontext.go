package svc

import (
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
	immodels.ConversationsModel
	immodels.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		ChatLogModel:       immodels.NewChatLogModel(c.Mongo.Url, c.Mongo.Db, "chat_log"),
		ConversationsModel: immodels.NewConversationsModel(c.Mongo.Url, c.Mongo.Db, "conversations"),
		ConversationModel:  immodels.NewConversationModel(c.Mongo.Url, c.Mongo.Db, "conversation"),
	}
}
