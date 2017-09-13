package esp4g

import (
	"google.golang.org/grpc"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

func NewGrpcServer(pb string, proxyAddress string, extensionAddress string, cfg config.ExtensionConfig) *grpc.Server {
	fds, err := utils.ReadFileDescriptorSet(pb)
	if err != nil {
		utils.Logger.Fatalw("failed to read pb file", "pb", pb, "err", err)
	}

	services := createProxyServiceDesc(fds)

	logInterceptor := newAccessLogInterceptor(extensionAddress, fds, cfg)
	controlInterceptor := newAccessControlInterceptor(extensionAddress, fds, cfg)

	opts := []grpc.ServerOption{}

	{
		f := controlInterceptor.createApiKeyInterceptor(nil)
		g := logInterceptor.createAccessLogInterceptor(f)
		opts = append(opts, grpc.UnaryInterceptor(*g))
	}

	{
		f := controlInterceptor.createStreamApiKeyInterceptor(nil)
		g := logInterceptor.createStreamAccessLogInterceptor(f)
		opts = append(opts, grpc.StreamInterceptor(*g))
	}

	server := grpc.NewServer(opts...)

	proxy, err := newProxyServer(proxyAddress)
	if err != nil {
		utils.Logger.Fatalw("failed to create new proxy server", "err", err)
	}

	for _, desc := range services {
		server.RegisterService(&desc, proxy)
	}

	return server
}
