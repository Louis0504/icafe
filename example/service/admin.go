package service

import (
	"context"
	"github.com/ymetas/icafe/example/gen-go/proto/admin"
)

type AdminService struct {
	admin.UnimplementedAdminServiceServer
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (r AdminService) GetContent(ctx context.Context, in *admin.GetContentParam) (*admin.GetContentResponse, error) {
	return &admin.GetContentResponse{
		Name: "1",
	}, nil
}
