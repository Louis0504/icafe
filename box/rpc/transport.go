package rpc

import (
	"context"
	"net/http"
	"strings"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/sirupsen/logrus"
)

const (
	defaultBufferSize = 4096
)

func newTransport(ctx context.Context, client *http.Client, targetURL string, serviceName, method string, headers map[string]string) thrift.TTransport {
	customHeaders := http.Header{
		"X-ZONE-API":           []string{serviceName + "." + method},
		"X-ZONE-ORIGIN":        []string{"ZvideoService"},
		"X-ZONE-ORIGIN-APP":    []string{"zvideo"},
	}

	for k, v := range headers {
		if strings.HasPrefix(k, "X-ZONE") {
			continue
		}
		customHeaders[k] = []string{v}
	}

	transport, err := thrift.NewTHttpClientWithOptions(targetURL, thrift.THttpClientOptions{Client: client})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"target_url": targetURL,
			"error_msg":  err.Error(),
		}).Panicln("client: Provided url is malformed.")
	}

	httpClient := transport.(*thrift.THttpClient)
	for key, value := range customHeaders {
		httpClient.SetHeader(key, value[0])
	}

	return thrift.NewTBufferedTransport(transport, defaultBufferSize)
}
