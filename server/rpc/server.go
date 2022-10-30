package rpc

import (
	"net/http"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
)

func NewServer(services map[string]thrift.TProcessor) *Server {
	s := &Server{
		processor:       thrift.NewTMultiplexedProcessor(),
		protocolFactory: thrift.NewTBinaryProtocolFactoryDefault(),
	}
	for serviceName, serviceProcessor := range services {
		s.processor.RegisterProcessor(serviceName, serviceProcessor)
	}

	return s
}

type Server struct {
	httpServer      *http.Server
	processor       *thrift.TMultiplexedProcessor
	protocolFactory *thrift.TBinaryProtocolFactory
	middlewares     []func(http.Handler) http.Handler
}

func (s *Server) Use(middlewares ...func(http.Handler) http.Handler) *Server {
	s.middlewares = append(s.middlewares, middlewares...)
	return s
}

func (s *Server) Chain(h http.Handler) http.Handler {
	for i := range s.middlewares {
		h = s.middlewares[len(s.middlewares)-1-i](h)
	}
	return h
}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/check_health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("zhi~"))
	})
	//mux.Handle("/", s.Chain(http.HandlerFunc(s.Handler)))

	s.httpServer = &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}
