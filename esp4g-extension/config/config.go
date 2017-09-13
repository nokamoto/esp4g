package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ExtensionConfig struct {
	Logs Logs `yaml:"logs"`

	Authentication Authentication `yaml:"authentication"`

	Usage Usage `yaml:"usage"`
}

func FromYaml(yml []byte) (ExtensionConfig, error) {
	var cfg ExtensionConfig
	return cfg, yaml.Unmarshal(yml, &cfg)
}

func FromYamlFile(file string) (ExtensionConfig, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return ExtensionConfig{}, err
	}
	return FromYaml(buf)
}
