package rpc

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/ymetas/icafe/log"
	"github.com/ymetas/icafe/server"
	"net/http"
)

type TRPCBundle struct {
	name   string
	server *Server

	serviceMap       map[string]thrift.TProcessor
	listenAddr       string
	extraMiddlewares []func(http.Handler) http.Handler
}

func NewTRPCBundle(name string, opts ...TRPCOption) server.Bundle {
	defaults := getDefaults()
	s := &TRPCBundle{
		name:       name,
		listenAddr: defaults.listenAddr,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.server = NewServer(s.serviceMap)
	s.server.Use(s.extraMiddlewares...)

	return s
}

func (s *TRPCBundle) Type() string {
	return "tzone"
}

func (s *TRPCBundle) Name() string {
	return s.name
}

func (s *TRPCBundle) Run(ctx context.Context) error {
	return s.server.Run(s.listenAddr)
}

func (s *TRPCBundle) Stop() context.Context {
	ctx2, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := s.server.Close()
		if err != nil {
			log.Errorf("Close tzone service error: %v", err)
		}
		ctx2.Done()
	}()

	return ctx2
}
