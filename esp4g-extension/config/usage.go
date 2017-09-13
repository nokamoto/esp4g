package config

type Requirement struct {
	ProviderId string `yaml:"provider_id"`
}

type Rule struct {
	Selector string `yaml:"selector"`

	AllowUnregisteredCalls bool `yaml:"allow_unregistered_calls"`

	Requirements *[]Requirement `yaml:"requirements"`
}

type Usage struct {
	Rules []Rule `yaml:"rules"`
}
