package esp4g_test

import "testing"

func TestAllowRegisteredApiKeys(t *testing.T) {
	config := "access_control-allow.yaml"
	apiKey := "guest"
	checkUnaryProxy(t, config, apiKey)
	checkClientSideStream(t, config, apiKey)
	checkServerSideStream(t, config, apiKey)
	checkBidirectionalStream(t, config, apiKey)
}
