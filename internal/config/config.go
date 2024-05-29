package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	PORT            string
	GITHUB_TOKEN    string
	SERVER_BASE_URL string
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

	server_base_url, err := getEnv("SERVER_BASE_URL")
	if err != nil {
		return nil, err
	}

	if !strings.Contains(server_base_url, ":") {
		server_base_url = fmt.Sprintf("%s:%s", server_base_url, port)
	}

	config := &Config{
		PORT:            port,
		GITHUB_TOKEN:    github_token,
		SERVER_BASE_URL: server_base_url,
	}

	return config, nil
}
