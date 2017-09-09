package esp4g

import (
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
)

func NewGrpcServer(pb string, proxyPort int, accessLogPort int, accessControlPort int) *grpc.Server {
	data, err := ioutil.ReadFile(pb)
	if err != nil {
		Logger.Fatalw("failed to read pb file", "pb", pb, "err", err)
	}

	fds := &descriptor.FileDescriptorSet{}

	proto.Unmarshal(data, fds)

	services := createProxyServiceDesc(fds)

	logInterceptor := newAccessLogInterceptor(accessLogPort)
	controlInterceptor := newAccessControlInterceptor(accessControlPort)

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

	proxy, err := newProxyServer(proxyPort)
	if err != nil {
		Logger.Fatalw("failed to create new proxy server", "err", err)
	}

	for _, desc := range services {
		server.RegisterService(&desc, proxy)
	}

	return server
}
