package main

import (
	"flag"
	"net"
	"fmt"
	"github.com/nokamoto/esp4g/esp4g/esp4g"
)

func main() {
	var (
		pb      = flag.String("d", "descriptor.pb", "FileDescriptorSet protocol buffer file")
		port    = flag.Int("p", 9000, "The gRPC server port")
		proxy   = flag.Int("proxy", 8000, "The gRPC proxy port")
		log     = flag.Int("log", 10000, "The gRPC access log service port")
		control = flag.Int("control", 10000, "The gRPC access control service port")
	)

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		esp4g.Logger.Fatalw("failed to listen", "err", err)
	} else {
		esp4g.Logger.Infow("listen port", "port", port)
	}

	server := esp4g.NewGrpcServer(*pb, *proxy, *log, *control)

	esp4g.Logger.Infow("start esp server")
	server.Serve(lis)
}
