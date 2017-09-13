package esp4g_test

import (
	"testing"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

func checkDeny(t *testing.T, file string, keys... string) {
	yaml, err := config.FromYamlFile(file)
	if err != nil {
		t.Error(err)
	}

	checkDenyUnary(t, yaml, keys...)
	checkDenyCStream(t, yaml, keys...)
	checkDenySStream(t, yaml, keys...)
	checkDenyBStream(t, yaml, keys...)
}

func TestAllowRegisteredApiKeys(t *testing.T) {
	yaml, err := config.FromYamlFile("yaml/allow-keys.yaml")
	if err != nil {
		t.Error(err)
	}

	apiKey := "guest"

	checkUnary(t, yaml, apiKey)
	checkCStream(t, yaml, apiKey)
	checkSStream(t, yaml, apiKey)
	checkBStream(t, yaml, apiKey)
}

func TestDenyUnregisteredCalls(t *testing.T) {
	checkDeny(t, "yaml/deny.yaml")
	checkDeny(t, "yaml/allow-keys.yaml", "guest2")
}

func TestComplexAccessControl(t *testing.T) {
	yaml, err := config.FromYamlFile("yaml/complex.yaml")
	if err != nil {
		t.Error(err)
	}

	t.Log("allow_unregistered_calls: true")
	checkUnary(t, yaml)

	t.Log("allways denied")
	checkDenyCStream(t, yaml)
	checkDenyCStream(t, yaml, "guest")
	checkDenyCStream(t, yaml, "admin")

	t.Log("provider_id: guest")
	checkDenySStream(t, yaml)
	checkSStream(t, yaml, "guest")
	checkDenySStream(t, yaml, "admin")

	t.Log("provider_id: admin")
	checkDenyBStream(t, yaml)
	checkDenyBStream(t, yaml, "guest")
	checkBStream(t, yaml, "admin")
}
