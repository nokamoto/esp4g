package extension

import (
	"golang.org/x/net/context"
	proto "github.com/nokamoto/esp4g/protobuf"
)

type accessControlService struct {
	rules []Rule
}

func (a *accessControlService)Access(_ context.Context, id *proto.AccessIdentity) (*proto.AccessControl, error) {
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

func newAccessControlService(config Config) *accessControlService {
	return &accessControlService{rules: config.Usage.Rules}
}