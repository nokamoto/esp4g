package esp4g_test

import (
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
)

func newGrpcServer(port int) (*grpc.Server, *PingService) {
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	ps := &PingService{}
	ping.RegisterPingServiceServer(server, ps)

	return server, ps
}
