package pipeline

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Step struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

type Config struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
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
