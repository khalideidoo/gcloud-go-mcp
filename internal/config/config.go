// Package config provides configuration management for the gcloud MCP server.
package config

import (
	"os"
	"time"
)

// Config holds all configuration for the gcloud MCP server.
type Config struct {
	// Project is the default GCP project ID.
	Project string

	// Region is the default region for regional resources.
	Region string

	// Zone is the default zone for zonal resources.
	Zone string

	// GCloudPath is the path to the gcloud binary.
	GCloudPath string

	// CommandTimeout is the maximum duration for command execution.
	CommandTimeout time.Duration
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	return &Config{
		Project:        getEnv("GCLOUD_PROJECT", ""),
		Region:         getEnv("GCLOUD_REGION", ""),
		Zone:           getEnv("GCLOUD_ZONE", ""),
		GCloudPath:     getEnv("GCLOUD_PATH", "gcloud"),
		CommandTimeout: getDurationEnv("GCLOUD_TIMEOUT", 5*time.Minute),
	}
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getDurationEnv returns the value of an environment variable as a duration or a default value.
func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
