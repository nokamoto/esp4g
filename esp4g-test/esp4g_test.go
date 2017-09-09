package esp4g_test

import (
	"testing"
	"github.com/nokamoto/esp4g/esp4g/esp4g"
	"github.com/nokamoto/esp4g/esp4g-extension/esp4g-extension"
	"net"
	"fmt"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"golang.org/x/net/context"
)

const UNARY_DESCRIPTOR = "unary-descriptor.pb"
const CONFIG = "config.yaml"
const PROXY_PORT = 9000
const UPSTREAM_PORT = 8000
const EXTENSION_PORT = 10000

func TestUnaryProxy(t *testing.T) {
	proxyServer := esp4g.NewGrpcServer(UNARY_DESCRIPTOR, UPSTREAM_PORT, EXTENSION_PORT, EXTENSION_PORT)
	extensionServer := extension.NewGrpcServer(CONFIG)
	upstreamServer, service := newGrpcServer(UPSTREAM_PORT)

	proxy := make(chan error, 1)
	ext := make(chan error, 1)
	upstream := make(chan error, 1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", UPSTREAM_PORT))
		if err != nil {
			upstream <- err
		} else {
			upstream <- upstreamServer.Serve(lis)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", EXTENSION_PORT))
		if err != nil {
			ext <- err
		} else {
			ext <- extensionServer.Serve(lis)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", PROXY_PORT))
		if err != nil {
			proxy <- err
		} else {
			proxy <- proxyServer.Serve(lis)
		}
	}()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("localhost:%d", PROXY_PORT), opts...)
	if err != nil {
		t.Error(err)
	}

	defer con.Close()

	client := ping.NewPingServiceClient(con)

	req := &ping.Ping{X: 100}
	res, err := client.Send(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	if len(service.requests) != 1 || *req != service.requests[0] {
		t.Errorf("unexpected request: %v %v", req, service.requests)
	}
	if len(service.responses) != 1 || *res != service.responses[0] {
		t.Errorf("unexpected response: %v %v", res, service.responses)
	}

	upstreamServer.GracefulStop()
	extensionServer.GracefulStop()
	proxyServer.GracefulStop()
}


