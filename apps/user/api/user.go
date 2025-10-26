package main

import (
	"flag"
	"fmt"

	"easy-chat/apps/user/api/internal/config"
	"easy-chat/apps/user/api/internal/handler"
	"easy-chat/apps/user/api/internal/svc"
	"easy-chat/pkg/configserver"
	"easy-chat/pkg/resultx"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

// var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	// conf.MustLoad(*configFile, &c)
	var configs = "user-api.yaml"
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.18.48:3379",
		ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
		Namespace:      "user",
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
		// proc.WrapUp() //  停止服务
		// wg.Add(1)

		// go func(c config.Config) {
		// 	defer wg.Done()
		// 	Run(c)
		// }(c)

		return nil
	})

	if err != nil {
		panic(err)
	}

	Run(c)

	// wg.Add(1)
	// go func(c config.Config) {
	// 	defer wg.Done()
	// 	Run(c)
	// }(c)
	// wg.Wait()
}

func Run(c config.Config) {
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
