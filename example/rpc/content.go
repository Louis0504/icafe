package rpc

import (
	"context"
	"github.com/Louis0504/icafe/example/gen-go/thrift/content_thrift/content"
)

type ContentServiceRPC interface {
	GetContent(ctx context.Context, param content.GetContentParam) content.GetContentResponse
}
