package extension

import (
	"google.golang.org/grpc"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	proto "github.com/nokamoto/esp4g/protobuf"
)

func NewGrpcServer(yml string) *grpc.Server {
	buf, err := ioutil.ReadFile(yml)
	if err != nil {
		Logger.Fatalw("failed to read yaml", "yaml", yml, "err", err)
	}

	var config Config
	if err = yaml.Unmarshal(buf, &config); err != nil {
		Logger.Fatalw("failed to unmarshal", "err", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	proto.RegisterAccessControlServiceServer(server, newAccessControlService(config))
	proto.RegisterAccessLogServiceServer(server, newAccessLogService(config))

	return server
}
