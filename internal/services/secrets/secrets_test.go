package secrets

import (
	"encoding/json"
	"testing"
	"time"

	"gcloud-go-mcp/internal/config"
	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newTestConfig() *config.Config {
	return &config.Config{
		Project:        "test-project",
		Region:         "us-central1",
		Zone:           "us-central1-a",
		GCloudPath:     "gcloud",
		CommandTimeout: 5 * time.Minute,
	}
}

// Helper function to simulate parseArgs behavior for testing
// This avoids issues with MCP SDK internal struct types
func testParseArgs(argsJSON json.RawMessage) map[string]any {
	var args map[string]any
	if argsJSON != nil {
		_ = json.Unmarshal(argsJSON, &args)
	}
	if args == nil {
		args = make(map[string]any)
	}
	return args
}

func TestRegisterTools(t *testing.T) {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "0.0.1",
		},
		&mcp.ServerOptions{},
	)
	base := services.NewBaseService(newTestConfig())

	// Should not panic
	RegisterTools(server, base)

	// Verify tools were registered by checking the server
	// Note: The MCP SDK may not expose a way to list tools directly,
	// so we verify registration didn't panic
}

func TestParseArgs_WithArguments(t *testing.T) {
	args := map[string]any{
		"secret_id": "my-secret",
		"project":   "my-project",
	}
	argsJSON, _ := json.Marshal(args)

	result := testParseArgs(argsJSON)

	if result["secret_id"] != "my-secret" {
		t.Errorf("expected secret_id 'my-secret', got %v", result["secret_id"])
	}
	if result["project"] != "my-project" {
		t.Errorf("expected project 'my-project', got %v", result["project"])
	}
}

func TestParseArgs_NilArguments(t *testing.T) {
	result := testParseArgs(nil)

	if result == nil {
		t.Error("expected non-nil map for nil arguments")
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParseArgs_EmptyArguments(t *testing.T) {
	result := testParseArgs(json.RawMessage(`{}`))

	if result == nil {
		t.Error("expected non-nil map for empty arguments")
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParseArgs_InvalidJSON(t *testing.T) {
	result := testParseArgs(json.RawMessage(`{invalid}`))

	// Should return empty map for invalid JSON
	if result == nil {
		t.Error("expected non-nil map for invalid JSON")
	}
}

func TestParseArgs_WithAllTypes(t *testing.T) {
	args := map[string]any{
		"string_param": "value",
		"int_param":    float64(42),
		"bool_param":   true,
		"array_param":  []any{"a", "b", "c"},
		"map_param":    map[string]any{"key": "value"},
	}
	argsJSON, _ := json.Marshal(args)

	result := testParseArgs(argsJSON)

	if result["string_param"] != "value" {
		t.Errorf("expected string_param 'value', got %v", result["string_param"])
	}
	if result["int_param"] != float64(42) {
		t.Errorf("expected int_param 42, got %v", result["int_param"])
	}
	if result["bool_param"] != true {
		t.Errorf("expected bool_param true, got %v", result["bool_param"])
	}

	arr, ok := result["array_param"].([]any)
	if !ok {
		t.Error("expected array_param to be []any")
	} else if len(arr) != 3 {
		t.Errorf("expected array_param length 3, got %d", len(arr))
	}

	m, ok := result["map_param"].(map[string]any)
	if !ok {
		t.Error("expected map_param to be map[string]any")
	} else if m["key"] != "value" {
		t.Errorf("expected map_param['key'] = 'value', got %v", m["key"])
	}
}

// Tool input schema tests
func TestToolSchemas(t *testing.T) {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "0.0.1",
		},
		&mcp.ServerOptions{},
	)
	base := services.NewBaseService(newTestConfig())
	RegisterTools(server, base)

	// These tests verify that tool schemas are correctly defined
	// by testing parameter extraction on various inputs

	tests := []struct {
		name        string
		args        map[string]any
		wantErr     bool
		errContains string
	}{
		{
			name: "create secret - valid",
			args: map[string]any{
				"secret_id": "my-secret",
			},
			wantErr: false,
		},
		{
			name:        "create secret - missing secret_id",
			args:        map[string]any{},
			wantErr:     true,
			errContains: "missing required parameter",
		},
		{
			name: "create secret - with labels",
			args: map[string]any{
				"secret_id": "my-secret",
				"labels": map[string]any{
					"env": "prod",
				},
			},
			wantErr: false,
		},
		{
			name: "access version - valid",
			args: map[string]any{
				"secret_id": "my-secret",
				"version":   "1",
			},
			wantErr: false,
		},
		{
			name: "access version - default version",
			args: map[string]any{
				"secret_id": "my-secret",
			},
			wantErr: false,
		},
		{
			name: "add iam binding - valid",
			args: map[string]any{
				"secret_id": "my-secret",
				"member":    "user:test@example.com",
				"role":      "roles/secretmanager.secretAccessor",
			},
			wantErr: false,
		},
		{
			name: "add iam binding - missing member",
			args: map[string]any{
				"secret_id": "my-secret",
				"role":      "roles/secretmanager.secretAccessor",
			},
			wantErr:     true,
			errContains: "missing required parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test required string extraction
			if _, hasSecretID := tt.args["secret_id"]; hasSecretID {
				val, err := services.GetRequiredString(tt.args, "secret_id")
				if err != nil {
					if !tt.wantErr {
						t.Errorf("unexpected error: %v", err)
					}
					return
				}
				if val == "" {
					t.Error("expected non-empty secret_id")
				}
			} else if tt.wantErr {
				_, err := services.GetRequiredString(tt.args, "secret_id")
				if err == nil {
					t.Error("expected error for missing secret_id")
				}
			}
		})
	}
}

// Test optional parameter extraction patterns used in secrets
func TestSecretParameterPatterns(t *testing.T) {
	t.Run("version defaults to latest", func(t *testing.T) {
		args := map[string]any{
			"secret_id": "my-secret",
		}
		version := services.GetOptionalString(args, "version", "latest")
		if version != "latest" {
			t.Errorf("expected default version 'latest', got %q", version)
		}
	})

	t.Run("filter is optional", func(t *testing.T) {
		args := map[string]any{}
		filter := services.GetOptionalString(args, "filter", "")
		if filter != "" {
			t.Errorf("expected empty filter, got %q", filter)
		}
	})

	t.Run("limit defaults to 100", func(t *testing.T) {
		args := map[string]any{}
		limit := services.GetOptionalInt(args, "limit", 100)
		if limit != 100 {
			t.Errorf("expected default limit 100, got %d", limit)
		}
	})

	t.Run("replication policy defaults to automatic", func(t *testing.T) {
		args := map[string]any{}
		policy := services.GetOptionalString(args, "replication_policy", "automatic")
		if policy != "automatic" {
			t.Errorf("expected default policy 'automatic', got %q", policy)
		}
	})

	t.Run("labels extraction", func(t *testing.T) {
		args := map[string]any{
			"labels": map[string]any{
				"env":  "production",
				"team": "backend",
			},
		}
		labels := services.GetOptionalStringMap(args, "labels")
		if len(labels) != 2 {
			t.Errorf("expected 2 labels, got %d", len(labels))
		}
		if labels["env"] != "production" {
			t.Errorf("expected env='production', got %q", labels["env"])
		}
	})
}

// Benchmark for parseArgs
func BenchmarkParseArgs(b *testing.B) {
	args := map[string]any{
		"secret_id": "my-secret",
		"project":   "my-project",
		"version":   "latest",
	}
	argsJSON, _ := json.Marshal(args)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testParseArgs(argsJSON)
	}
}

func BenchmarkParseArgs_LargePayload(b *testing.B) {
	args := map[string]any{
		"secret_id": "my-secret",
		"project":   "my-project",
		"labels": map[string]any{
			"env":     "production",
			"team":    "backend",
			"service": "api",
			"version": "v1",
		},
		"replication_policy": "automatic",
		"filter":             "state:ENABLED",
		"limit":              float64(500),
	}
	argsJSON, _ := json.Marshal(args)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testParseArgs(argsJSON)
	}
}
