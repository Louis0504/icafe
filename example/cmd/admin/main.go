package main

import (
	"github.com/ymetas/icafe/example/gen-go/proto/admin"
	"github.com/ymetas/icafe/example/service"
	"github.com/ymetas/icafe/server"
	"github.com/ymetas/icafe/server/grpc"
)

func main() {
	app := server.NewApplication(server.SentryIncludePaths("github.com/alonegrowing/admin-service"))

	grpcBundle := grpc.NewGRPCBundle("grpc-server", grpc.GRPCListen(":9999"))
	admin.RegisterAdminServiceServer(grpcBundle.Server, service.NewAdminService())

	app.AddBundle(grpcBundle)

	app.Run()
}
