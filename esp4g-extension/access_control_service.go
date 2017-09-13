package extension

import (
	"golang.org/x/net/context"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

type AccessControlService struct {
	rules []config.Rule
}

func (a *AccessControlService)Access(_ context.Context, id *proto.AccessIdentity) (*proto.AccessControl, error) {
	allow := proto.AccessControl{Policy: proto.AccessPolicy_ALLOW}
	for _, rule := range a.rules {
		if rule.Selector == id.Method {
			if rule.AllowUnregisteredCalls {
				return &allow, nil
			}
			for _, key := range id.ApiKey {
				for _, registeredKey := range rule.RegisteredApiKey {
					if key == registeredKey {
						return &allow, nil
					}
				}
			}
		}
	}
	return &proto.AccessControl{Policy: proto.AccessPolicy_DENY}, nil
}

func NewAccessControlService(cfg config.ExtensionConfig, _ *descriptor.FileDescriptorSet) *AccessControlService {
	return &AccessControlService{rules: cfg.Usage.Rules}
}
