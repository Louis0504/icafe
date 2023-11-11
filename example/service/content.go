package service

import (
	"context"
	"github.com/YLeseclaireurs/icafe/example/gen-go/thrift/content_thrift/base"
	"github.com/YLeseclaireurs/icafe/example/gen-go/thrift/content_thrift/content"
)

type ContentService struct {
}

func NewContentService() *ContentService {
	return &ContentService{}
}

func (r ContentService) GetContent(ctx context.Context, in *content.GetContentParam) (*content.GetContentResponse, error) {
	return &content.GetContentResponse{Content: &base.Content{
		ID:        1,
		Name:      "levin",
		Content:   "一条内容",
		UpdatedAt: "2020-04-02 12:00:00",
		CreatedAt: "2020-04-02 13:00:00",
	}}, nil
}
