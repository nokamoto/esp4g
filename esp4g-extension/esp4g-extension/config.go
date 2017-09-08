package extension

type Rule struct {
	Selector string `yaml:"selector"`

	AllowUnregisteredCalls bool `yaml:"allow_unregistered_calls"`

	RegisteredApiKey []string `yaml:"registered_api_keys"`
}

type Usage struct {
	Rules []Rule `yaml:"rules"`
}

type Config struct {
	Usage Usage `yaml:"usage"`
}
