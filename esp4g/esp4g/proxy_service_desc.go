package esp4g

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	"fmt"
	"golang.org/x/net/context"
	"github.com/nokamoto/esp4g/esp4g-utils"
)

func createProxyHandler(method string) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := newProxyMessage()
		if err := dec(in); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return srv.(proxyServer).Proxy(ctx, method, in)
		}

		info := &grpc.UnaryServerInfo{
			Server: srv,
			FullMethod: method,
		}

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(proxyServer).Proxy(ctx, method, in)
		}

		return interceptor(ctx, in, info, handler)
	}
}

func createClientSideStreamingProxyHandler(method string, desc *grpc.StreamDesc) func(interface{}, grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(proxyServer).ProxyClientSideStreaming(method, &clientSideServerStream{stream}, desc)
	}
}

func createServerSideStreamingProxyHandler(method string, desc *grpc.StreamDesc) func(interface{}, grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		m := new(proxyMessage)
		if err := stream.RecvMsg(m); err != nil {
			return err
		}
		return srv.(proxyServer).ProxyServerSideStreaming(method, m, &serverSideServerStream{stream}, desc)
	}
}

func createBidirectionalStreamingProxyHandler(method string, desc *grpc.StreamDesc) func(interface{}, grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return srv.(proxyServer).ProxyBidirectionalStreaming(method, &bidirectionalServerStream{stream}, desc)
	}
}

func nilOrBool(x *bool) bool {
	return x != nil && *x
}

func createServiceDesc(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) grpc.ServiceDesc {
	utils.Logger.Infow("service descriptor", "service", service)

	serviceName := fmt.Sprintf("%s.%s", file.GetPackage(), service.GetName())

	methods := make([]grpc.MethodDesc, 0)
	streams := make([]grpc.StreamDesc, 0)

	for _, method := range service.Method {
		methodName := utils.GetFullMethodName(file, service, method)

		cs := nilOrBool(method.ClientStreaming)
		ss := nilOrBool(method.ServerStreaming)

		if cs || ss {
			desc := grpc.StreamDesc{
				StreamName: method.GetName(),
				ServerStreams: ss,
				ClientStreams: cs,
			}

			streams = append(streams, desc)
			last := len(streams) - 1

			if cs && ss {
				streams[last].Handler = createBidirectionalStreamingProxyHandler(methodName, &streams[last])
			} else if cs {
				streams[last].Handler = createClientSideStreamingProxyHandler(methodName, &streams[last])
			} else if ss {
				streams[last].Handler = createServerSideStreamingProxyHandler(methodName, &streams[last])
			}
		} else {
			handler := createProxyHandler(methodName)

			desc := grpc.MethodDesc{
				MethodName: method.GetName(),
				Handler: handler,
			}

			methods = append(methods, desc)
		}
	}

	return grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*proxyServer)(nil),
		Metadata: file.GetName(),
		Methods: methods,
		Streams: streams,
	}
}

func createProxyServiceDesc(fds *descriptor.FileDescriptorSet) []grpc.ServiceDesc {
	services := make([]grpc.ServiceDesc, 0)

	for _, file := range fds.GetFile() {
		utils.Logger.Infow("file descriptor", "file", file.GetName())

		for _, service := range file.GetService() {
			services = append(services, createServiceDesc(file, service))
		}
	}

	return services
}
