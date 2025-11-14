package config

import "os"

// Config holds application configuration
type Config struct {
	SonarURL   string
	SonarToken string
}

// New creates a new configuration from environment variables
func New() *Config {
	return &Config{
		SonarURL:   os.Getenv("SONAR_URL"),
		SonarToken: os.Getenv("SONAR_TOKEN"),
	}
}
