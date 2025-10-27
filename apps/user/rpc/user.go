package main

import (
	"flag"
	"fmt"

	"easy-chat/apps/user/rpc/internal/config"
	"easy-chat/apps/user/rpc/internal/server"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/configserver"
	"easy-chat/pkg/interceptor/rpcserver"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// conf.MustLoad(*configFile, &c)

	// ctx := svc.NewServiceContext(c)

	// if err := ctx.SetRootToken(); err != nil {
	// 	panic(err)
	// }

	// s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
	// 	user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

	// 	if c.Mode == service.DevMode || c.Mode == service.TestMode {
	// 		reflection.Register(grpcServer)
	// 	}
	// })
	// s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	// defer s.Stop()

	// fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	// s.Start()

	var configs = "user-rpc.yaml"
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.18.48:3379",
		ProjectKey:     "dca19ffa163054feef33432fad5f9833",
		Namespace:      "user",
		Configs:        configs,
		ConfigFilePath: "../etc/conf",
		// 本地测试使用以下配置
		// ConfigFilePath: "./etc/conf",
		LogLevel: "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		err := configserver.LoadFromJsonBytes(bytes, &c)
		if err != nil {
			fmt.Println("config read err :", err)
			return nil
		}
		fmt.Println(configs, "config has changed")
		return nil
	})
	if err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)

	if err := ctx.SetRootToken(); err != nil {
		panic(err)
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddUnaryInterceptors(rpcserver.LogInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
