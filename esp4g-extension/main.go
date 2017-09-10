package main

import (
	"flag"
	"net"
	"fmt"
	"github.com/nokamoto/esp4g/esp4g-extension/esp4g-extension"
	"github.com/nokamoto/esp4g/esp4g-utils"
)

func main() {
	var (
		port = flag.Int("p", 10000, "The gRPC server port")
		yml = flag.String("c", "./config.yaml", "The application config file path")
		pb = flag.String("d", "./descriptor.pb", "FileDescriptorSet protocol buffer file")
	)

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		utils.Logger.Fatalw("failed to listen", "err", err)
	} else {
		utils.Logger.Infow("listen port", "port", *port)
	}

	server := extension.NewGrpcServer(*yml, *pb)

	utils.Logger.Info("start esp extension server")
	server.Serve(lis)
}
