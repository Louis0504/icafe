package rpc

import (
	"github.com/apache/thrift/lib/go/thrift"
	"net/http"
)

type TRPCOption func(*TRPCBundle)

func TRPCListen(listenAddr string) TRPCOption {
	return func(s *TRPCBundle) {
		s.listenAddr = listenAddr
	}
}

func WithTRPCServiceMap(serviceMap map[string]thrift.TProcessor) TRPCOption {
	return func(s *TRPCBundle) {
		s.serviceMap = serviceMap
	}
}

func WithMiddlewares(middlewares ...func(http.Handler) http.Handler) TRPCOption {
	return func(s *TRPCBundle) {
		s.extraMiddlewares = append(s.extraMiddlewares, middlewares...)
	}
}
