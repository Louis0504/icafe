package grpc



type GRPCOption func(bundle *GRPCBundle)

func GRPCListen(listenAddr string) GRPCOption {
	return func(s *GRPCBundle) {
		s.listenAddr = listenAddr
	}
}
