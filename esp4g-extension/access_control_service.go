package main

import (
	"golang.org/x/net/context"
	extension "github.com/nokamoto/esp4g/protobuf"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type accessControlService struct {
}

func (*accessControlService)Access(context.Context, *extension.AccessIdentity) (*extension.AccessControl, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

func NewAccessControlService() *accessControlService {
	return &accessControlService{}
}
