package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"time"
)

type Server = grpc.Server
type ServerOption = grpc.ServerOption

func NewServer(opts ...ServerOption) *Server {
	overwriteOpts := []ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      5 * time.Minute,
			MaxConnectionAgeGrace: 10 * time.Second,
			Time:                  time.Second,
			Timeout:               time.Millisecond * 100,
		}),
		//grpc.ChainUnaryInterceptor(tgrpc.WithServerTelemetry),
	}
	opts = append(opts, overwriteOpts...)

	s := grpc.NewServer(opts...)
	reflection.Register(s)
	return s
}
