package extension

import (
	"golang.org/x/net/context"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type AccessControlService struct {
	auth *authenticationProviders
}

func (a *AccessControlService)Access(_ context.Context, id *proto.AccessIdentity) (*proto.AccessControl, error) {
	policy, err := a.auth.Allow(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.AccessControl{Policy: policy}, nil
}

func NewAccessControlService(cfg config.ExtensionConfig, fds *descriptor.FileDescriptorSet) *AccessControlService {
	methods := utils.Methods(fds)
	auth, err := NewAuthenticationProviders(cfg.Authentication, cfg.Usage.Rules, methods)
	if err != nil {
		utils.Logger.Fatalw("failed to build authentication m", "err", err)
	}
	return &AccessControlService{auth: auth}
}
