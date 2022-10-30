package grpc

import (
	"context"
	"net"

	"github.com/pkg/errors"
)

type GRPCBundle struct {
	name       string
	Server     *Server
	listenAddr string
}

func NewGRPCBundle(name string, opts ...GRPCOption) *GRPCBundle {
	defaults := getDefaults()
	s := &GRPCBundle{
		name:       name,
		listenAddr: defaults.listenAddr,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.Server = NewServer()

	return s
}

func (s *GRPCBundle) Type() string {
	return "gRPC"
}

func (s *GRPCBundle) Name() string {
	return s.name
}

func (s *GRPCBundle) Run(ctx context.Context) error {
	addr, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return errors.Wrap(err, "listen failed")
	}
	return s.Server.Serve(addr)
}

func (s *GRPCBundle) Stop() context.Context {
	ctx2, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		s.Server.Stop()
		ctx2.Done()
	}()

	return ctx2
}
