package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigFilePath string

type Config struct {
	Log      Log      `yaml:"log"`
	Auth     Auth     `yaml:"auth"`
	Database Database `yaml:"database"`
}

func NewConfig(filePath ConfigFilePath) (Config, error) {
	configBytes, err := os.ReadFile(string(filePath))
	if err != nil {
		return Config{}, fmt.Errorf("failed to read YAML file: %w", err)
	}

	config := Config{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return config, nil
}
