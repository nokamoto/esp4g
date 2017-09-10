package extension

import (
	"google.golang.org/grpc"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/nokamoto/esp4g/esp4g-utils"
)

func NewGrpcServer(yml string, pb string) *grpc.Server {
	fds, err := utils.ReadFileDescriptorSet(pb)
	if err != nil {
		utils.Logger.Fatalw("failed to read pb file", "pb", pb, "err", err)
	}

	buf, err := ioutil.ReadFile(yml)
	if err != nil {
		utils.Logger.Fatalw("failed to read yaml", "yaml", yml, "err", err)
	}

	var config Config
	if err = yaml.Unmarshal(buf, &config); err != nil {
		utils.Logger.Fatalw("failed to unmarshal", "err", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	proto.RegisterAccessControlServiceServer(server, newAccessControlService(config, fds))
	proto.RegisterAccessLogServiceServer(server, newAccessLogService(config, fds))

	return server
}
