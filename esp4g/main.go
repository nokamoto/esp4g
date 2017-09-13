package main

import (
	"flag"
	"net"
	"fmt"
	"github.com/nokamoto/esp4g/esp4g/esp4g"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

func main() {
	var (
		pb        = flag.String("d", "./descriptor.pb", "FileDescriptorSet protocol buffer file")
		port      = flag.Int("p", 9000, "The gRPC server port")
		yaml      = flag.String("c", "./config.yaml", "The application config file")
		extension = flag.String("e", "", "The gRPC extension service address (default in-process)")
		proxy     = flag.String("u", "localhost:8000", "The gRPC upstream address")
	)

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		utils.Logger.Fatalw("failed to listen", "err", err)
	} else {
		utils.Logger.Infow("listen port", "port", port)
	}

	cfg, err := config.FromYamlFile(*yaml)
	if err != nil {
		utils.Logger.Fatalw("failed to read yaml", "yaml", yaml, "err", err)
	}

	server := esp4g.NewGrpcServer(*pb, *proxy, *extension, cfg)

	utils.Logger.Infow("start esp server")
	server.Serve(lis)
}
