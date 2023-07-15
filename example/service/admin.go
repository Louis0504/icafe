package service

import (
	"context"
	"github.com/Louis0504/icafe/example/gen-go/proto/admin"
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
