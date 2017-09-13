package extension

import (
	"google.golang.org/grpc"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
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

	var cfg config.ExtensionConfig
	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		utils.Logger.Fatalw("failed to unmarshal", "err", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	proto.RegisterAccessControlServiceServer(server, NewAccessControlService(cfg, fds))
	proto.RegisterAccessLogServiceServer(server, NewAccessLogService(cfg, fds))

	return server
}
