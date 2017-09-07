package main

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	"log"
	"fmt"
	"golang.org/x/net/context"
)

func createProxyHandler(method string) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := NewProxyMessage()
		if err := dec(in); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return srv.(ProxyServer).Proxy(ctx, method, in)
		}

		info := &grpc.UnaryServerInfo{
			Server: srv,
			FullMethod: method,
		}

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProxyServer).Proxy(ctx, method, in)
		}

		return interceptor(ctx, in, info, handler)
	}
}

func createServiceDesc(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) grpc.ServiceDesc {
	log.Printf("service: %v", service)

	serviceName := fmt.Sprintf("%s.%s", file.GetPackage(), service.GetName())

	methods := make([]grpc.MethodDesc, 0)
	streams := make([]grpc.StreamDesc, 0)

	for _, method := range service.Method {
		methodName := fmt.Sprintf("/%s/%s", serviceName, method.GetName())
		handler := createProxyHandler(methodName)

		desc := grpc.MethodDesc{
			MethodName: method.GetName(),
			Handler: handler,
		}

		methods = append(methods, desc)
	}

	return grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*ProxyServer)(nil),
		Metadata: file.GetName(),
		Methods: methods,
		Streams: streams,
	}
}

func CreateProxyServiceDesc(fds *descriptor.FileDescriptorSet) []grpc.ServiceDesc {
	services := make([]grpc.ServiceDesc, 0)

	for _, file := range fds.GetFile() {
		log.Printf("file: %v", file.GetName())

		for _, service := range file.GetService() {
			services = append(services, createServiceDesc(file, service))
		}
	}

	return services
}
