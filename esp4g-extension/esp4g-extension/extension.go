package extension

import (
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	proto "github.com/nokamoto/esp4g/protobuf"
)

func NewGrpcServer(yml string) *grpc.Server {
	buf, err := ioutil.ReadFile(yml)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err = yaml.Unmarshal(buf, &config); err != nil {
		log.Fatal(err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	proto.RegisterAccessControlServiceServer(server, NewAccessControlService(config))
	proto.RegisterAccessLogServiceServer(server, NewAccessLogService(config))

	return server
}
