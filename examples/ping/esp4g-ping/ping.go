package main

import (
	"flag"
	"fmt"
	"os"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"time"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/ptypes/empty"
)

var (
	host = flag.String("h", "localhost", "The gRPC server host")
	port = flag.Int("p", 9000, "The gRPC server port")
	interval = flag.Int("n", 1, "Wait n seconds")
	unavailable = flag.Bool("u", false, "Send ping to Unavailable")
	apikey = flag.String("k", "guest", "The gRPC request 'x-api-key'")
	check = flag.Bool("test", false, "Health check only")
)

type PerRPCCredentials struct {}

func (PerRPCCredentials)GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"x-api-key": *apikey}, nil
}

func (PerRPCCredentials)RequireTransportSecurity() bool {
	return false
}

func send(client ping.PingServiceClient, req ping.Ping) (*ping.Pong, error) {
	if *unavailable {
		return client.Unavailable(context.Background(), &req)
	}
	return client.Send(context.Background(), &req)
}

func main() {
	flag.Parse()

	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithPerRPCCredentials(PerRPCCredentials{})}

	con, err := grpc.Dial(fmt.Sprintf("%s:%d", *host, *port), opts...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer con.Close()

	if *check {
		_, err := ping.NewHealthCheckServiceClient(con).Check(context.Background(), &empty.Empty{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("ok")
		}
	} else {
		client := ping.NewPingServiceClient(con)

		count := int64(0)
		for count >= 0 {
			req := ping.Ping{X: count}

			if res, err := send(client, req); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(res)
			}

			time.Sleep(time.Duration(*interval) * time.Second)

			count = count + 1
		}
	}
}
