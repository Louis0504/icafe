package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
)


type Client struct {
	targetName  string
	timeout     time.Duration
	serviceName string
	discovery   *Discovery
	hostPort    string
	headers     map[string]string
	client      *http.Client
}


type Option func(*Client)

// TargetName targetName use service discovery to find an available address for specified name
func TargetName(name string) Option {
	return Option(func(c *Client) {
		c.targetName = name
		c.discovery = NewDiscovery(name)
	})
}

// Timeout timeout specify the timeout for the underline transport.
//
// default value is 500ms.
func Timeout(t time.Duration) Option {
	return func(c *Client) {
		c.timeout = t
	}
}

// HostPort use host and port for remote service provider.
//
// if both targetName and HostPort provided, HostPort is used, targetName is ignored.
func HostPort(host string, port string) Option {
	return func(c *Client) {
		c.hostPort = host + ":" + port
	}
}

func Url(url string) Option {
	return func(c *Client) {
		c.hostPort = url
	}
}

// Headers Custom http headers
func Headers(headers map[string]string) Option {
	return func(c *Client) {
		c.headers = headers
	}
}

func (c *Client) SetHeader(key string, value string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}

	c.headers[key] = value
}

func (t *Client) Call(ctx context.Context, method string, args, result thrift.TStruct) (thrift.ResponseMeta,error) {
	meta := thrift.ResponseMeta{}

	// Fetch address
	var url string
	var currentAddr *Address
	if t.hostPort != "" {
		url = "http://" + t.hostPort
	} else {
		var err error
		currentAddr, err = t.discovery.GetAddress()
		if err != nil {
			return meta, fmt.Errorf("GetAddress failed for %s. Error:%s\n", t.targetName, err)
		}
		url = "http://" + currentAddr.String()
	}

	// Make transport
	transport := newTransport(ctx, t.client, url, t.serviceName, method, t.headers)

	// Make protocol
	protocol := newProtocol(transport, t.serviceName)

	// Make real request
	conn := thrift.NewTStandardClient(protocol, protocol)
	var err error
	meta,err = conn.Call(ctx, method, args, result);
	if  err != nil {
		if _, ok := err.(thrift.TTransportException); ok && currentAddr != nil {
			t.discovery.DiscardAddress(currentAddr)
		}
	}

	_ = protocol.Transport().Close()
	return meta,err
}

// New create a new tzone client with specified service and options.
func New(serviceName string, opts ...Option) *Client {
	c := &Client{
		timeout:     500 * time.Millisecond,
		serviceName: serviceName,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.discovery == nil && c.hostPort == "" {
		panic("client: either targetName or HostPort option must be specified.")
	}

	// fork from https://github.com/golang/go/blob/release-branch.go1.11/src/net/http/transport.go#L42
	c.client = &http.Client{
		Timeout: c.timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          10240,
			MaxIdleConnsPerHost:   1024,
			IdleConnTimeout:       1 * time.Second, // 过长的存活时间链接可能会被 server 关掉
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return c
}
