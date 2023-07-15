package rpc

import (
	"context"

	"github.com/Louis0504/icafe/example/gen-go/proto/admin"
)

type AdminRPC interface {
	GetContent(ctx context.Context, param admin.GetContentParam) admin.GetContentResponse
}
