package main

import (
	"flag"
	"net"
	"fmt"
	"github.com/nokamoto/esp4g/esp4g/esp4g"
	"github.com/nokamoto/esp4g/esp4g-utils"
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
		utils.Logger.Fatalw("failed to listen", "err", err)
	} else {
		utils.Logger.Infow("listen port", "port", port)
	}

	server := esp4g.NewGrpcServer(*pb, *proxy, *log, *control)

	utils.Logger.Infow("start esp server")
	server.Serve(lis)
}
