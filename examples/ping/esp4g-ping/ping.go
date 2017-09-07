package main

import (
	"flag"
	"fmt"
	"os"
	"google.golang.org/grpc"
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"context"
	"time"
)

var (
	host = flag.String("h", "localhost", "The gRPC server host")
	port = flag.Int("p", 9000, "The gRPC server port")
	interval = flag.Int("n", 1, "Wait n seconds")
	unavailable = flag.Bool("u", false, "Send ping to Unavailable")
)

func send(client ping.PingServiceClient, req ping.Ping) (*ping.Pong, error) {
	if *unavailable {
		return client.Unavailable(context.Background(), &req)
	}
	return client.Send(context.Background(), &req)
}

func main() {
	flag.Parse()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("%s:%d", *host, *port), opts...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer con.Close()

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
