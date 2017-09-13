package extension

import (
	"google.golang.org/grpc"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

func NewGrpcServer(cfg config.ExtensionConfig, pb string) *grpc.Server {
	fds, err := utils.ReadFileDescriptorSet(pb)
	if err != nil {
		utils.Logger.Fatalw("failed to read pb file", "pb", pb, "err", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	proto.RegisterAccessControlServiceServer(server, NewAccessControlService(cfg, fds))
	proto.RegisterAccessLogServiceServer(server, NewAccessLogService(cfg, fds))

	return server
}
