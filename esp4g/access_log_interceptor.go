package main

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"time"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func createAccessLogInterceptor(log func(string, time.Duration, codes.Code, int, int)) grpc.UnaryServerInterceptor {
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		res, err := handler(ctx, req)

		elapsed := time.Since(start)

		code := codes.Unknown
		if stat, ok := status.FromError(err); ok {
			code = stat.Code()
		}

		inBytes := -1
		outBytes := -1
		if m, ok := req.(*ProxyMessage); ok {
			inBytes = len(m.bytes)
		}
		if m, ok := res.(*ProxyMessage); ok {
			outBytes = len(m.bytes)
		}

		log(info.FullMethod, elapsed, code, inBytes, outBytes)

		return res, err
	}
	return grpc.UnaryServerInterceptor(f)
}

func createStreamAccessLogInterceptor(log func(string, time.Duration, codes.Code)) grpc.StreamServerInterceptor {
	f := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		err := handler(srv, ss)

		elapsed := time.Since(start)

		code := codes.Unknown
		if stat, ok := status.FromError(err); ok {
			code = stat.Code()
		}

		log(info.FullMethod, elapsed, code)

		return err
	}
	return grpc.StreamServerInterceptor(f)
}
