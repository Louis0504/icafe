package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

type ClientConn = grpc.ClientConn
type UnaryInvoker = grpc.UnaryInvoker
type CallOption = grpc.CallOption
type DialOption = grpc.DialOption

func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
	defaultOpts := []grpc.DialOption{grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: 500 * time.Millisecond,
	})}
	opts = append(defaultOpts, opts...)

	overwriteOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(withClientTelemetry),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	}
	opts = append(opts, overwriteOpts...)

	return grpc.DialContext(ctx, target, opts...)
}

func withClientTelemetry(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error {
	return nil
}
