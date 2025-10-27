package main

import (
	"flag"
	"fmt"

	"easy-chat/apps/im/api/internal/config"
	"easy-chat/apps/im/api/internal/handler"
	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/pkg/configserver"
	"easy-chat/pkg/resultx"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// conf.MustLoad(*configFile, &c)
	var configs = "im-api.yaml"
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

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
