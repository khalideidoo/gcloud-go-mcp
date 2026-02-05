// Package services provides base functionality for GCP service implementations.
package services

import (
	"context"
	"fmt"

	"gcloud-go-mcp/internal/config"
	"gcloud-go-mcp/internal/executor"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BaseService provides common functionality for all GCP services.
type BaseService struct {
	Executor *executor.Executor
	Config   *config.Config
}

// NewBaseService creates a new base service.
func NewBaseService(cfg *config.Config) *BaseService {
	return &BaseService{
		Executor: executor.New(cfg),
		Config:   cfg,
	}
}

// ToolResult creates a successful tool result with text content.
func ToolResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

// ToolError creates an error tool result.
func ToolError(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: err.Error()},
		},
	}
}

// ToolHandler is the function signature for tool handlers.
type ToolHandler func(ctx context.Context, args map[string]any) (*mcp.CallToolResult, error)

// GetRequiredString extracts a required string parameter.
func GetRequiredString(args map[string]any, key string) (string, error) {
	val, ok := args[key]
	if !ok {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("parameter %s must be a string", key)
	}
	if str == "" {
		return "", fmt.Errorf("parameter %s cannot be empty", key)
	}
	return str, nil
}

// GetOptionalString extracts an optional string parameter.
func GetOptionalString(args map[string]any, key string, defaultVal string) string {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}
	str, ok := val.(string)
	if !ok {
		return defaultVal
	}
	return str
}

// GetOptionalInt extracts an optional integer parameter.
func GetOptionalInt(args map[string]any, key string, defaultVal int) int {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}
	// JSON numbers are float64
	num, ok := val.(float64)
	if !ok {
		return defaultVal
	}
	return int(num)
}

// GetOptionalBool extracts an optional boolean parameter.
func GetOptionalBool(args map[string]any, key string, defaultVal bool) bool {
	val, ok := args[key]
	if !ok {
		return defaultVal
	}
	b, ok := val.(bool)
	if !ok {
		return defaultVal
	}
	return b
}

// GetOptionalStringArray extracts an optional string array parameter.
func GetOptionalStringArray(args map[string]any, key string) []string {
	val, ok := args[key]
	if !ok {
		return nil
	}
	arr, ok := val.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// GetOptionalStringMap extracts an optional string map parameter.
func GetOptionalStringMap(args map[string]any, key string) map[string]string {
	val, ok := args[key]
	if !ok {
		return nil
	}
	m, ok := val.(map[string]any)
	if !ok {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		if s, ok := v.(string); ok {
			result[k] = s
		} else {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result
}
