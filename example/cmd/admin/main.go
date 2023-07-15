package main

import (
	"github.com/Louis0504/icafe/example/gen-go/proto/admin"
	"github.com/Louis0504/icafe/example/service"
	"github.com/Louis0504/icafe/server"
	"github.com/Louis0504/icafe/server/grpc"
)

func main() {
	app := server.NewApplication(server.SentryIncludePaths("github.com/alonegrowing/admin-service"))

	grpcBundle := grpc.NewGRPCBundle("grpc-server", grpc.GRPCListen(":9999"))
	admin.RegisterAdminServiceServer(grpcBundle.Server, service.NewAdminService())

	app.AddBundle(grpcBundle)

	app.Run()
}
