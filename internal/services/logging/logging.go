// Package logging provides MCP tools for Google Cloud Logging.
package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Cloud Logging tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// Read logs
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_logging_read",
			Description: "Read log entries with optional filtering",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"filter": map[string]any{
						"type":        "string",
						"description": "Log filter expression (e.g., 'resource.type=cloud_run_revision AND severity>=ERROR')",
					},
					"resource_type": map[string]any{
						"type":        "string",
						"description": "Resource type (e.g., cloud_run_revision, gce_instance, cloud_function)",
					},
					"log_name": map[string]any{
						"type":        "string",
						"description": "Specific log name to read from",
					},
					"severity": map[string]any{
						"type":        "string",
						"description": "Minimum severity level",
						"enum":        []string{"DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"},
					},
					"limit": map[string]any{
						"type":        "number",
						"description": "Maximum number of entries to return",
						"default":     50,
					},
					"freshness": map[string]any{
						"type":        "string",
						"description": "How far back to read (e.g., 1h, 30m, 1d)",
						"default":     "1h",
					},
					"order": map[string]any{
						"type":        "string",
						"description": "Sort order: asc or desc",
						"default":     "desc",
						"enum":        []string{"asc", "desc"},
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)

			// Build filter parts
			var filterParts []string

			if filter := services.GetOptionalString(args, "filter", ""); filter != "" {
				filterParts = append(filterParts, filter)
			}
			if resourceType := services.GetOptionalString(args, "resource_type", ""); resourceType != "" {
				filterParts = append(filterParts, fmt.Sprintf("resource.type=%s", resourceType))
			}
			if logName := services.GetOptionalString(args, "log_name", ""); logName != "" {
				filterParts = append(filterParts, fmt.Sprintf("logName:%s", logName))
			}
			if severity := services.GetOptionalString(args, "severity", ""); severity != "" {
				filterParts = append(filterParts, fmt.Sprintf("severity>=%s", severity))
			}

			cmd := base.Executor.Command("logging", "read")

			// Add filter as positional argument if present
			if len(filterParts) > 0 {
				filterStr := strings.Join(filterParts, " AND ")
				cmd = base.Executor.Command("logging", "read", filterStr)
			}

			cmd.WithProject(services.GetOptionalString(args, "project", "")).
				WithFlag("limit", fmt.Sprintf("%d", services.GetOptionalInt(args, "limit", 50))).
				WithFlag("freshness", services.GetOptionalString(args, "freshness", "1h"))

			order := services.GetOptionalString(args, "order", "desc")
			if order == "asc" {
				cmd.WithFlag("order", "asc")
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// List logs
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_logging_logs_list",
			Description: "List available logs in a project",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)

			result, err := base.Executor.Command("logging", "logs", "list").
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Write log
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_logging_write",
			Description: "Write a log entry",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"log_name", "message"},
				"properties": map[string]any{
					"log_name": map[string]any{
						"type":        "string",
						"description": "Log name to write to",
					},
					"message": map[string]any{
						"type":        "string",
						"description": "Log message",
					},
					"severity": map[string]any{
						"type":        "string",
						"description": "Log severity",
						"default":     "INFO",
						"enum":        []string{"DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRITICAL"},
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			logName, err := services.GetRequiredString(args, "log_name")
			if err != nil {
				return services.ToolError(err), nil
			}
			message, err := services.GetRequiredString(args, "message")
			if err != nil {
				return services.ToolError(err), nil
			}

			// Create JSON payload
			payload := fmt.Sprintf(`{"message":"%s"}`, message)

			result, err := base.Executor.Command("logging", "write", logName, payload).
				WithFlag("severity", services.GetOptionalString(args, "severity", "INFO")).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithTextFormat().
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			if result.Stdout == "" {
				return services.ToolResult("Log entry written successfully"), nil
			}
			return services.ToolResult(result.Stdout), nil
		},
	)
}

func parseArgs(req *mcp.CallToolRequest) map[string]any {
	var args map[string]any
	if req.Params.Arguments != nil {
		_ = json.Unmarshal(req.Params.Arguments, &args)
	}
	if args == nil {
		args = make(map[string]any)
	}
	return args
}
