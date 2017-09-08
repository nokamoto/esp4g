package main

type Server struct {
	Port int `yaml:"port"`
}

type Rule struct {
	Selector string `yaml:"selector"`

	AllowUnregisteredCalls bool `yaml:"allow_unregistered_calls"`

	RegisteredApiKey []string `yaml:"registered_api_keys"`
}

type Usage struct {
	Rules []Rule `yaml:"rules"`
}

type Config struct {
	Server Server `yaml:"server"`

	Usage Usage `yaml:"usage"`
}
