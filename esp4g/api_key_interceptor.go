package main

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type Policy int

const (
	ALLOW Policy = iota
	DENY
)

const API_KEY_HEADER = "x-api-key"

func createApiKeyInterceptor(allow func(string, []string) Policy, next *grpc.UnaryServerInterceptor) *grpc.UnaryServerInterceptor {
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "failed to read metadata")
		}

		apiKey, ok := md[API_KEY_HEADER]

		if allow(info.FullMethod, apiKey) == ALLOW {
			if next != nil {
				return (*next)(ctx, req, info, handler)
			}
			return handler(ctx, req)
		}

		return nil, status.Error(codes.Unauthenticated, "unauthenticated request")
	}

	i := grpc.UnaryServerInterceptor(f)

	return &i
}

func createStreamApiKeyInterceptor(allow func(string, []string) Policy, next *grpc.StreamServerInterceptor) *grpc.StreamServerInterceptor {
	f := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Error(codes.Internal, "failed to read metadata")
		}

		apiKey, ok := md[API_KEY_HEADER]

		if allow(info.FullMethod, apiKey) == ALLOW {
			if next != nil {
				return (*next)(srv, ss, info, handler)
			}
			return handler(srv, ss)
		}

		return status.Error(codes.Unauthenticated, "unauthenticated request")
	}

	i := grpc.StreamServerInterceptor(f)

	return &i
}
