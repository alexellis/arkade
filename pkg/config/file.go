package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ArkadeConfig struct {
	Ignore []string `yaml:"ignore"`
}

func Load(file string) (*ArkadeConfig, error) {

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &ArkadeConfig{}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
