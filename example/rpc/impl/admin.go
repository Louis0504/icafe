package impl

import (
	"context"
	"github.com/ymetas/icafe/example/gen-go/proto/admin"
	"github.com/ymetas/icafe/server/grpc"
)

type AdminRPCImpl struct {
	adminClient *admin.AdminServiceClient
}

var DefaultAdminRPC *AdminRPCImpl

func init() {
	DefaultAdminRPC = NewAdminRPC()
}

func NewAdminRPC() *AdminRPCImpl {
	client, _ := grpc.DialContext(context.Background(), "127.0.0.1:9000")
	adminClient := admin.NewAdminServiceClient(client)
	return &AdminRPCImpl{
		adminClient: &adminClient,
	}
}

func (r *AdminRPCImpl) GetContent(param admin.GetContentParam) admin.GetContentResponse {
	return admin.GetContentResponse{}
}
