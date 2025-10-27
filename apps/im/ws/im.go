package main

import (
	"easy-chat/apps/im/ws/internal/config"
	"easy-chat/apps/im/ws/internal/handler"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/pkg/configserver"
	"flag"
	"fmt"
	"time"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// conf.MustLoad(*configFile, &c)
	var configs = "im-ws.yaml"
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.18.48:3379",
		ProjectKey:     "dca19ffa163054feef33432fad5f9833",
		Namespace:      "im",
		Configs:        configs,
		ConfigFilePath: "../etc/conf",
		// ConfigFilePath: "./etc/conf",
		LogLevel: "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		err := configserver.LoadFromJsonBytes(bytes, &c)
		if err != nil {
			fmt.Println("config read err :", err)
		}
		fmt.Println(configs, "config has changed")

		return nil
	})

	if err != nil {
		panic(err)
	}

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerMaxConnectionIdle(1000*time.Second),
		websocket.WithServerAck(websocket.NoAck),
	)
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server at ", c.ListenOn)
	srv.Start()
}
