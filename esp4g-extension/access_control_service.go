package main

import (
	"golang.org/x/net/context"
	extension "github.com/nokamoto/esp4g/protobuf"
)

type accessControlService struct {
	rules []Rule
}

func (a *accessControlService)Access(_ context.Context, id *extension.AccessIdentity) (*extension.AccessControl, error) {
	allow := extension.AccessControl{Policy: extension.AccessPolicy_ALLOW}
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
	return &extension.AccessControl{Policy: extension.AccessPolicy_DENY}, nil
}

func NewAccessControlService(config Config) *accessControlService {
	return &accessControlService{rules: config.Usage.Rules}
}
