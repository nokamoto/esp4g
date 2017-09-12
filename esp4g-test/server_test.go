package esp4g_test

import (
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"fmt"
	"github.com/nokamoto/esp4g/esp4g/esp4g"
	"net"
	"testing"
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"time"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/nokamoto/esp4g/esp4g-extension"
)

const UNARY_DESCRIPTOR = "unary-descriptor.pb"
const STREAM_DESCRIPTOR = "stream-descriptor.pb"
const PROXY_PORT = 9000
const UPSTREAM_PORT = 8000
const EXTENSION_PORT = 10000

func newGrpcServer() (*grpc.Server, *PingService, *CalcService) {
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	ps := &PingService{}
	ping.RegisterPingServiceServer(server, ps)

	cs := &CalcService{}
	calc.RegisterCalcServiceServer(server, cs)

	return server, ps, cs
}

func preflightPing(t *testing.T, con *grpc.ClientConn) {
	i := 0
	for i < 10 {
		_, err := callPing(con, &ping.Ping{})

		if err == nil {
			return
		}

		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unauthenticated {
			return
		}

		t.Logf("%d: %v", i, err)

		time.Sleep(time.Duration(i * 100) * time.Millisecond)

		i = i + 1
	}
	t.Error("preflight timed out")
}

func preflightCalc(t *testing.T, con *grpc.ClientConn) {
	i := 0
	for i < 10 {
		_, err := callCalcCStream(con, []*calc.Operand{{}})

		if err == nil {
			return
		}

		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unauthenticated {
			t.Logf("got unauthenticated error: %s", stat.Message())
			return
		}

		t.Logf("%d: %v", i, err)

		time.Sleep(time.Duration(i * 100) * time.Millisecond)

		i = i + 1
	}
	t.Error("preflight timed out")
}

func inproc(t *testing.T, descriptor string, config string, c chan error) (*grpc.Server, *grpc.Server) {
	proxy := esp4g.NewGrpcServer(
		descriptor,
		fmt.Sprintf("localhost:%d", UPSTREAM_PORT),
		"",
		config,
	)

	start(t, proxy, PROXY_PORT, c)

	return proxy, nil
}

func outproc(t *testing.T, descriptor string, config string, cp chan error, se chan error) (*grpc.Server, *grpc.Server) {
	proxy := esp4g.NewGrpcServer(
		descriptor,
		fmt.Sprintf("localhost:%d", UPSTREAM_PORT),
		fmt.Sprintf("localhost:%d", EXTENSION_PORT),
		"",
	)

	start(t, proxy, PROXY_PORT, cp)

	ext := extension.NewGrpcServer(config, descriptor)

	start(t, ext, EXTENSION_PORT, se)

	return proxy, ext
}

func start(t *testing.T, server *grpc.Server, port int, c chan error) {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			t.Logf("listen error: %v", err)
			c <- err
		} else {
			t.Logf("listen %d", port)
			c <- server.Serve(lis)
		}
	}()
}

func stop(t *testing.T, server *grpc.Server, c chan error) {
	server.GracefulStop()

	select {
	case err := <- c:
		t.Log("gracefully shutdown", err)
	}
}

func run(t *testing.T, f func(*grpc.ClientConn, *PingService, *CalcService)) {
	upstream := make(chan error, 1)
	upstreamServer, ps, cs := newGrpcServer()
	start(t, upstreamServer, UPSTREAM_PORT, upstream)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("localhost:%d", PROXY_PORT), opts...)
	if err != nil {
		t.Error(err)
	}

	defer con.Close()

	f(con, ps, cs)

	stop(t, upstreamServer, upstream)
}

func withServers(t *testing.T, descriptor string, config string, f func(*grpc.ClientConn, *PingService, *CalcService)) {
	proxy := make(chan error, 1)
	p, e := inproc(t, descriptor, config, proxy)

	t.Log("run inproc")
	run(t, f)

	stop(t, p, proxy)

	proxy = make(chan error, 1)
	ext := make(chan error, 1)
	p, e = outproc(t, descriptor, config, proxy, ext)

	t.Log("run outproc")
	run(t, f)

	stop(t, p, proxy)
	stop(t, e, ext)
}
