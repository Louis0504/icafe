package rpc

import "github.com/apache/thrift/lib/go/thrift"

func newProtocol(transport thrift.TTransport, serviceName string) thrift.TProtocol {
	protocol := thrift.NewTBinaryProtocolTransport(transport)
	return thrift.NewTMultiplexedProtocol(protocol, serviceName)
}
