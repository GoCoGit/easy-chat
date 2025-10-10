package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	sync.RWMutex

	authentication Authentication

	routes   map[string]HandlerFunc
	addr     string
	upgrader websocket.Upgrader

	patten string
	opt    *serverOption

	connToUser map[*Conn]string
	userToConn map[string]*Conn

	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		upgrader: websocket.Upgrader{},

		patten:         opt.patten,
		authentication: opt.Authentication,
		opt:            &opt,

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		Logger: logx.WithContext(context.Background()),
	}
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			s.Errorf("ServerWs panic: %v", err)
		}
	}()

	// conn, err := s.upgrader.Upgrade(w, r, nil)
	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}
	// if err != nil {
	// 	s.Errorf("Upgrade error: %v", err)
	// 	return
	// }

	if !s.authentication.Auth(w, r) {
		// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存访问权限")))
		s.Send(&Message{FrameType: FrameData, Data: "不存访问权限"}, conn)
		conn.Close()
		return
	}

	// 记录链接
	s.AddConn(conn, r)

	// 处理链接
	go s.handleConn(conn)
}

func (s *Server) AddConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 不允许重复登录，如果已经登录，则关闭旧的连接
	if c := s.userToConn[uid]; c != nil {
		fmt.Println("close old conn")
		s.Close(c)
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

// 根据连接对象进行任务处理
func (s *Server) handleConn(conn *Conn) {
	uids := s.GetUsers(conn)
	conn.Uid = uids[0]

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("ReadMessage error: %v", err)
			s.Close(conn)
			return
		}

		// 解析消息
		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			s.Errorf("Unmarshal error: %v", err)
			s.Close(conn)
			return
		}

		// 根据消息类型进行处理
		switch message.FrameType {
		case FramePing:
			s.Send(&Message{FrameType: FramePing}, conn)
		case FrameData:
			// 根据请求method分发路由
			if handler, ok := s.routes[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)}, conn)
			}
		}
	}
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("停止服务")
}
