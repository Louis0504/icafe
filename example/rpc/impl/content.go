package impl

import (
	"context"

	"github.com/Louis0504/icafe/example/gen-go/thrift/content_thrift/content"
	client "github.com/Louis0504/icafe/server/rpc"
	"time"
)

type ContentRPCImpl struct {
	contentClient *content.ContentServiceClient
}

var DefaultContentRPC *ContentRPCImpl

func init() {
	DefaultContentRPC = NewHermesRPC()
}

func NewHermesRPC() *ContentRPCImpl {
	tClient := client.New(
		"ContentService",
		client.TargetName("content-thrift"),
		client.Timeout(200*time.Millisecond),
	)
	return &ContentRPCImpl{
		contentClient: content.NewContentServiceClient(tClient),
	}
}

func (r *ContentRPCImpl) GetContent(ctx context.Context, in content.GetContentParam) content.GetContentResponse {
	return content.GetContentResponse{}
}
