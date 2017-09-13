package config

type Provider struct {
	Id string `yaml:"id"`

	RegisteredApiKeys []string `yaml:"registered_api_keys"`
}

type Authentication struct {
	Providers []Provider `yaml:"providers"`
}
