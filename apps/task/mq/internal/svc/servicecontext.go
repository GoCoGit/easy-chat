package svc

import (
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/task/mq/internal/config"
	"easy-chat/pkg/constants"
	"net/http"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	config.Config

	WsClient websocket.Client
	*redis.Redis
	immodels.ChatLogModel
	immodels.ConversationModel

	socialclient.Social // 增加社交RPC，用于获取群组成员
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      immodels.NewChatLogModel(c.Mongo.Url, c.Mongo.Db, "chat_log"),
		ConversationModel: immodels.NewConversationModel(c.Mongo.Url, c.Mongo.Db, "conversation"),

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}

	token, err := svc.GetSystemRootToken()
	if err != nil {
		panic(err)
	}
	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))

	return svc
}

func (svc *ServiceContext) GetSystemRootToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
