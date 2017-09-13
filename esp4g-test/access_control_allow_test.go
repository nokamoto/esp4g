package esp4g_test

import (
	"testing"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

func TestAllowRegisteredApiKeys(t *testing.T) {
	yaml, err := config.FromYamlFile("yaml/allow-keys.yaml")
	if err != nil {
		t.Error(err)
	}

	apiKey := "guest"

	t.Log(yaml)

	checkUnaryProxy(t, yaml, apiKey)
	checkClientSideStream(t, yaml, apiKey)
	checkServerSideStream(t, yaml, apiKey)
	checkBidirectionalStream(t, yaml, apiKey)
}
