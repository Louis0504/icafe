package rpc

import (
	"context"

	"github.com/ymetas/icafe/example/gen-go/proto/admin"
)

type AdminRPC interface {
	GetContent(ctx context.Context, param admin.GetContentParam) admin.GetContentResponse
}
