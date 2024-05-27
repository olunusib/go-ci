package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	GITHUB_TOKEN string
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return value, nil
}

func Load() (*Config, error) {
	port, err := getEnv("PORT")
	if err != nil {
		port = "8080"
	}

	github_token, err := getEnv("GITHUB_TOKEN")
	if err != nil {
		return nil, err
	}

	config := &Config{
		Port:         port,
		GITHUB_TOKEN: github_token,
	}

	return config, nil
}
