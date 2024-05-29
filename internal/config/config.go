package config

import (
	"fmt"
	"net/url"
	"os"
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

	parsedURL, err := url.Parse(server_base_url)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_BASE_URL: %w", err)
	}

	if parsedURL.Port() == "" {
		parsedURL.Host = fmt.Sprintf("%s:%s", parsedURL.Hostname(), port)
		server_base_url = parsedURL.String()
	}

	config := &Config{
		PORT:            port,
		GITHUB_TOKEN:    github_token,
		SERVER_BASE_URL: server_base_url,
	}

	return config, nil
}
