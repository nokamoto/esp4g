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
	"golang.org/x/net/context"
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
		client := ping.NewPingServiceClient(con)

		_, err := client.Send(context.Background(), &ping.Ping{})

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
		client := calc.NewCalcServiceClient(con)

		stream, err := client.AddAll(context.Background())
		if err != nil {
			t.Error(err)
		}

		_, err = stream.CloseAndRecv()

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

func inproc(descriptor string, config string) (*grpc.Server, *grpc.Server) {
	proxy := esp4g.NewGrpcServer(
		descriptor,
		fmt.Sprintf("localhost:%d", UPSTREAM_PORT),
		"",
		config,
	)

	start(proxy, PROXY_PORT)

	return proxy, nil
}

func outproc(descriptor string, config string) (*grpc.Server, *grpc.Server) {
	proxy := esp4g.NewGrpcServer(
		descriptor,
		fmt.Sprintf("localhost:%d", UPSTREAM_PORT),
		fmt.Sprintf("localhost:%d", EXTENSION_PORT),
		"",
	)

	start(proxy, PROXY_PORT)

	ext := extension.NewGrpcServer(config, descriptor)

	start(ext, EXTENSION_PORT)

	return proxy, ext
}

func start(server *grpc.Server, port int) {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			server.Serve(lis)
		}
	}()
}

func run(servers []*grpc.Server, t *testing.T, f func(*grpc.ClientConn, *PingService, *CalcService)) {
	upstreamServer, ps, cs := newGrpcServer()
	start(upstreamServer, UPSTREAM_PORT)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("localhost:%d", PROXY_PORT), opts...)
	if err != nil {
		t.Error(err)
	}

	defer con.Close()

	f(con, ps, cs)

	for _, server := range servers {
		if server != nil {
			server.GracefulStop()
		}
	}

	upstreamServer.GracefulStop()
}

func withServers(t *testing.T, descriptor string, config string, f func(*grpc.ClientConn, *PingService, *CalcService)) {
	p, e := inproc(descriptor, config)

	t.Log("run inproc")
	run([]*grpc.Server{p, e}, t, f)

	p, e = outproc(descriptor, config)

	t.Log("run outproc")
	run([]*grpc.Server{p, e}, t, f)
}
