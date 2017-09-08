package main

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"time"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func createAccessLogInterceptor(log func(string, time.Duration, codes.Code, int, int), next *grpc.UnaryServerInterceptor) *grpc.UnaryServerInterceptor {
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		var res interface{}
		var err error

		if next != nil {
			res, err = (*next)(ctx, req, info, handler)
		} else {
			res, err = handler(ctx, req)
		}

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

	i := grpc.UnaryServerInterceptor(f)

	return &i
}

func createStreamAccessLogInterceptor(log func(string, time.Duration, codes.Code), next *grpc.StreamServerInterceptor) *grpc.StreamServerInterceptor {
	f := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		var err error
		if next != nil {
			err = (*next)(srv, ss, info, handler)
		} else {
			err = handler(srv, ss)
		}

		elapsed := time.Since(start)

		code := codes.Unknown
		if stat, ok := status.FromError(err); ok {
			code = stat.Code()
		}

		log(info.FullMethod, elapsed, code)

		return err
	}

	i := grpc.StreamServerInterceptor(f)

	return &i
}
