package main

import (
	"golang.org/x/net/context"
	extension "github.com/nokamoto/esp4g/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type accessLogService struct {
}

func (*accessLogService)UnaryAccess(context.Context, *extension.UnaryAccessLog) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

func (*accessLogService)StreamAccess(context.Context, *extension.StreamAccessLog) (*empty.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

func NewAccessLogService() *accessLogService {
	return &accessLogService{}
}
