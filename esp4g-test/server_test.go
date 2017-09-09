package esp4g_test

import (
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"fmt"
	"github.com/nokamoto/esp4g/esp4g/esp4g"
	"github.com/nokamoto/esp4g/esp4g-extension/esp4g-extension"
	"net"
	"testing"
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"time"
)

func newGrpcServer(port int) (*grpc.Server, *PingService, *CalcService) {
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	ps := &PingService{}
	ping.RegisterPingServiceServer(server, ps)

	cs := &CalcService{}
	calc.RegisterCalcServiceServer(server, cs)

	return server, ps, cs
}

func withServers(t *testing.T, descriptor string, config string, f func(*grpc.ClientConn, *PingService, *CalcService)) {
	proxyServer := esp4g.NewGrpcServer(descriptor, UPSTREAM_PORT, EXTENSION_PORT, EXTENSION_PORT)
	extensionServer := extension.NewGrpcServer(config)
	upstreamServer, ps, cs := newGrpcServer(UPSTREAM_PORT)

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

	time.Sleep(1 * time.Second)

	f(con, ps, cs)

	upstreamServer.GracefulStop()
	extensionServer.GracefulStop()
	proxyServer.GracefulStop()
}
