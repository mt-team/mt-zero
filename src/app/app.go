package main

import (
	"flag"
	"fmt"

	"mtzero/src/app/app"
	"mtzero/src/app/internal/config"
	"mtzero/src/app/internal/server"
	"mtzero/src/app/internal/svc"
	bizResponse "mtzero/src/util/response"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/app.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewAppServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		app.RegisterAppServer(grpcServer, srv)
	})
	s.AddUnaryInterceptors(
		bizResponse.RecoveryInterceptor,
	)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
