package rpc

import (
	"context"
	"github.com/ymetas/icafe/example/gen-go/thrift/content_thrift/content"
)

type ContentServiceRPC interface {
	GetContent(ctx context.Context, param content.GetContentParam) content.GetContentResponse
}
