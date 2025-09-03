package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Shell      string `yaml:"shell"`
	RunCommand string `yaml:"run_command"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(path string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}
