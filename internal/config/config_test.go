package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("GCLOUD_PROJECT")
	os.Unsetenv("GCLOUD_REGION")
	os.Unsetenv("GCLOUD_ZONE")
	os.Unsetenv("GCLOUD_PATH")
	os.Unsetenv("GCLOUD_TIMEOUT")

	cfg := LoadConfig()

	if cfg.Project != "" {
		t.Errorf("expected empty Project, got %q", cfg.Project)
	}
	if cfg.Region != "" {
		t.Errorf("expected empty Region, got %q", cfg.Region)
	}
	if cfg.Zone != "" {
		t.Errorf("expected empty Zone, got %q", cfg.Zone)
	}
	if cfg.GCloudPath != "gcloud" {
		t.Errorf("expected GCloudPath 'gcloud', got %q", cfg.GCloudPath)
	}
	if cfg.CommandTimeout != 5*time.Minute {
		t.Errorf("expected CommandTimeout 5m, got %v", cfg.CommandTimeout)
	}
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("GCLOUD_PROJECT", "test-project")
	os.Setenv("GCLOUD_REGION", "us-west1")
	os.Setenv("GCLOUD_ZONE", "us-west1-a")
	os.Setenv("GCLOUD_PATH", "/custom/path/gcloud")
	os.Setenv("GCLOUD_TIMEOUT", "10m")

	defer func() {
		os.Unsetenv("GCLOUD_PROJECT")
		os.Unsetenv("GCLOUD_REGION")
		os.Unsetenv("GCLOUD_ZONE")
		os.Unsetenv("GCLOUD_PATH")
		os.Unsetenv("GCLOUD_TIMEOUT")
	}()

	cfg := LoadConfig()

	if cfg.Project != "test-project" {
		t.Errorf("expected Project 'test-project', got %q", cfg.Project)
	}
	if cfg.Region != "us-west1" {
		t.Errorf("expected Region 'us-west1', got %q", cfg.Region)
	}
	if cfg.Zone != "us-west1-a" {
		t.Errorf("expected Zone 'us-west1-a', got %q", cfg.Zone)
	}
	if cfg.GCloudPath != "/custom/path/gcloud" {
		t.Errorf("expected GCloudPath '/custom/path/gcloud', got %q", cfg.GCloudPath)
	}
	if cfg.CommandTimeout != 10*time.Minute {
		t.Errorf("expected CommandTimeout 10m, got %v", cfg.CommandTimeout)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		envValue   string
		defaultVal string
		want       string
	}{
		{
			name:       "returns env value when set",
			key:        "TEST_VAR_1",
			envValue:   "custom-value",
			defaultVal: "default",
			want:       "custom-value",
		},
		{
			name:       "returns default when not set",
			key:        "TEST_VAR_2",
			envValue:   "",
			defaultVal: "default",
			want:       "default",
		},
		{
			name:       "returns empty string when set to empty",
			key:        "TEST_VAR_3",
			envValue:   "",
			defaultVal: "default",
			want:       "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestGetDurationEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		envValue   string
		defaultVal time.Duration
		want       time.Duration
	}{
		{
			name:       "parses valid duration",
			key:        "TEST_DUR_1",
			envValue:   "30s",
			defaultVal: time.Minute,
			want:       30 * time.Second,
		},
		{
			name:       "parses minutes",
			key:        "TEST_DUR_2",
			envValue:   "15m",
			defaultVal: time.Minute,
			want:       15 * time.Minute,
		},
		{
			name:       "returns default for invalid duration",
			key:        "TEST_DUR_3",
			envValue:   "invalid",
			defaultVal: 2 * time.Minute,
			want:       2 * time.Minute,
		},
		{
			name:       "returns default when not set",
			key:        "TEST_DUR_4",
			envValue:   "",
			defaultVal: 3 * time.Minute,
			want:       3 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := getDurationEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getDurationEnv(%q, %v) = %v, want %v", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}
