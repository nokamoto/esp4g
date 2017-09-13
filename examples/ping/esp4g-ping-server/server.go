package main

import (
	"flag"
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/ptypes/empty"
)

type PingServer struct {}

func (PingServer)Send(ctx context.Context, req *ping.Ping) (*ping.Pong, error) {
	log.Printf("ping=%v", req)
	return &ping.Pong{Y: req.GetX()}, nil
}

func (PingServer)Unavailable(context.Context, *ping.Ping) (*ping.Pong, error) {
	return nil, status.Error(codes.Unavailable, "this method always returns unavailable")
}

func (PingServer)Check(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

var (
	port = flag.Int("p", 8000, "The gRPC server port")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("listen %v port", *port)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	ping.RegisterPingServiceServer(server, PingServer{})
	ping.RegisterHealthCheckServiceServer(server, PingServer{})

	log.Println("start esp4g-ping-server...")
	server.Serve(lis)
}
