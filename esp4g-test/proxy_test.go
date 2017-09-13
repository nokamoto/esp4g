package esp4g_test

import (
	"testing"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

func TestProxy(t *testing.T) {
	yaml, err := config.FromYamlFile("yaml/allow.yaml")
	if err != nil {
		t.Error(err)
	}
	checkUnary(t, yaml)
	checkCStream(t, yaml)
	checkSStream(t, yaml)
	checkBStream(t, yaml)
}

