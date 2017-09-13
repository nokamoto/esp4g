package extension

import (
	"github.com/nokamoto/esp4g/esp4g-extension/config"
	proto "github.com/nokamoto/esp4g/protobuf"
	"errors"
	"fmt"
)

type authenticationProvider interface {
	Allow(id *proto.AccessIdentity) (proto.AccessPolicy, error)
}

type allowUnregisteredCallsProvider struct {}

func (*allowUnregisteredCallsProvider)Allow(id *proto.AccessIdentity) (proto.AccessPolicy, error) {
	return proto.AccessPolicy_ALLOW, nil
}

type localAuthenticationProvider struct {
	provider config.Provider
}

func newAuthenticationProvider(provider config.Provider) authenticationProvider {
	return &localAuthenticationProvider{provider: provider}
}

func (a *localAuthenticationProvider)Allow(id *proto.AccessIdentity) (proto.AccessPolicy, error) {
	for _, registered := range a.provider.RegisteredApiKeys {
		for _, key := range id.ApiKey {
			if key == registered {
				return proto.AccessPolicy_ALLOW, nil
			}
		}
	}
	return proto.AccessPolicy_DENY, nil
}

type authenticationProviders struct {
	m map[string][]authenticationProvider
}

func (a *authenticationProviders)Allow(id *proto.AccessIdentity) (proto.AccessPolicy, error) {
	if requirements, ok := a.m[id.Method]; !ok {
		return proto.AccessPolicy_DENY, nil
	} else if len(requirements) > 0 {
		for _, requirement := range requirements {
			if policy, err := requirement.Allow(id); err != nil || policy == proto.AccessPolicy_DENY {
				return proto.AccessPolicy_DENY, err
			}
		}
		return proto.AccessPolicy_ALLOW, nil
	}
	return proto.AccessPolicy_DENY, nil
}

func selectorMatch(selector string, method string) bool {
	return selector == method
}

func NewAuthenticationProviders(auth config.Authentication, rules []config.Rule, methods []string) (*authenticationProviders, error) {
	m := map[string][]authenticationProvider{}

	providers := map[string]authenticationProvider{}
	for _, provider := range auth.Providers {
		providers[provider.Id] = newAuthenticationProvider(provider)
	}

	for _, method := range methods {
		for _, rule := range rules {
			if selectorMatch(rule.Selector, method) {
				if rule.AllowUnregisteredCalls {
					m[method] = []authenticationProvider{&allowUnregisteredCallsProvider{}}
				} else if rule.Requirements != nil {
					requirements := []authenticationProvider{}

					for _, requirement := range *rule.Requirements {
						if provider, ok := providers[requirement.ProviderId]; !ok {
							s := fmt.Sprintf("provider-id undefined: rule=%v, requirement=%v", rule, requirement)
							return nil, errors.New(s)
						} else {
							requirements = append(requirements, provider)
						}
					}

					m[method] = requirements
				}
			}
		}
	}

	return &authenticationProviders{m: m}, nil
}
