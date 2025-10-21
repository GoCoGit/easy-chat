package handler

import (
	"easy-chat/apps/im/ws/internal/handler/conversation"
	"easy-chat/apps/im/ws/internal/handler/push"
	"easy-chat/apps/im/ws/internal/handler/user"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	// 注册路由，Method是模式，Handler是处理的回调
	// 将list传入server.go的AddRoutes方法，添加到server的routes中
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
		{
			Method:  "conversation.markChat",
			Handler: conversation.MarkRead(svc),
		},
	})
}
