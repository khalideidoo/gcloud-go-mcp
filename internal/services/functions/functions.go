// Package functions provides MCP tools for Google Cloud Functions.
package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Cloud Functions tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List functions
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_functions_list",
			Description: "List Cloud Functions",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)

			cmd := base.Executor.Command("functions", "list").
				WithProject(services.GetOptionalString(args, "project", ""))

			if region := services.GetOptionalString(args, "region", ""); region != "" {
				cmd.WithFlag("regions", region)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Describe function
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_functions_describe",
			Description: "Get details of a Cloud Function",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"function", "region"},
				"properties": map[string]any{
					"function": map[string]any{
						"type":        "string",
						"description": "Function name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region",
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
			function, err := services.GetRequiredString(args, "function")
			if err != nil {
				return services.ToolError(err), nil
			}
			region, err := services.GetRequiredString(args, "region")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("functions", "describe", function).
				WithRegion(region).
				WithProject(services.GetOptionalString(args, "project", "")).
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Deploy function
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_functions_deploy",
			Description: "Deploy a Cloud Function",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"function", "runtime", "region"},
				"properties": map[string]any{
					"function": map[string]any{
						"type":        "string",
						"description": "Function name",
					},
					"runtime": map[string]any{
						"type":        "string",
						"description": "Runtime (e.g., python311, nodejs18, go121)",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region",
					},
					"trigger_http": map[string]any{
						"type":        "boolean",
						"description": "Deploy as HTTP-triggered function",
					},
					"trigger_topic": map[string]any{
						"type":        "string",
						"description": "Pub/Sub topic to trigger on",
					},
					"trigger_bucket": map[string]any{
						"type":        "string",
						"description": "Cloud Storage bucket to trigger on",
					},
					"entry_point": map[string]any{
						"type":        "string",
						"description": "Function entry point",
					},
					"source": map[string]any{
						"type":        "string",
						"description": "Source code location (local path or GCS URL)",
					},
					"memory": map[string]any{
						"type":        "string",
						"description": "Memory limit (e.g., 256MB, 512MB)",
					},
					"timeout": map[string]any{
						"type":        "string",
						"description": "Timeout in seconds",
					},
					"env_vars": map[string]any{
						"type":        "object",
						"description": "Environment variables",
					},
					"service_account": map[string]any{
						"type":        "string",
						"description": "Service account email",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"allow_unauthenticated": map[string]any{
						"type":        "boolean",
						"description": "Allow unauthenticated invocations",
					},
					"gen2": map[string]any{
						"type":        "boolean",
						"description": "Deploy as 2nd generation function",
						"default":     true,
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			function, err := services.GetRequiredString(args, "function")
			if err != nil {
				return services.ToolError(err), nil
			}
			runtime, err := services.GetRequiredString(args, "runtime")
			if err != nil {
				return services.ToolError(err), nil
			}
			region, err := services.GetRequiredString(args, "region")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("functions", "deploy", function).
				WithFlag("runtime", runtime).
				WithRegion(region).
				WithProject(services.GetOptionalString(args, "project", ""))

			if services.GetOptionalBool(args, "gen2", true) {
				cmd.WithBoolFlag("gen2")
			}

			if services.GetOptionalBool(args, "trigger_http", false) {
				cmd.WithBoolFlag("trigger-http")
			}
			if topic := services.GetOptionalString(args, "trigger_topic", ""); topic != "" {
				cmd.WithFlag("trigger-topic", topic)
			}
			if bucket := services.GetOptionalString(args, "trigger_bucket", ""); bucket != "" {
				cmd.WithFlag("trigger-bucket", bucket)
			}
			if entryPoint := services.GetOptionalString(args, "entry_point", ""); entryPoint != "" {
				cmd.WithFlag("entry-point", entryPoint)
			}
			if source := services.GetOptionalString(args, "source", ""); source != "" {
				cmd.WithFlag("source", source)
			}
			if memory := services.GetOptionalString(args, "memory", ""); memory != "" {
				cmd.WithFlag("memory", memory)
			}
			if timeout := services.GetOptionalString(args, "timeout", ""); timeout != "" {
				cmd.WithFlag("timeout", timeout)
			}
			if sa := services.GetOptionalString(args, "service_account", ""); sa != "" {
				cmd.WithFlag("service-account", sa)
			}

			if envVars := services.GetOptionalStringMap(args, "env_vars"); len(envVars) > 0 {
				var pairs []string
				for k, v := range envVars {
					pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
				}
				cmd.WithFlag("set-env-vars", strings.Join(pairs, ","))
			}

			if services.GetOptionalBool(args, "allow_unauthenticated", false) {
				cmd.WithBoolFlag("allow-unauthenticated")
			}

			result, err := cmd.ExecuteWithRegion(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Delete function
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_functions_delete",
			Description: "Delete a Cloud Function",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"function", "region"},
				"properties": map[string]any{
					"function": map[string]any{
						"type":        "string",
						"description": "Function name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region",
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
			function, err := services.GetRequiredString(args, "function")
			if err != nil {
				return services.ToolError(err), nil
			}
			region, err := services.GetRequiredString(args, "region")
			if err != nil {
				return services.ToolError(err), nil
			}

			_, err = base.Executor.Command("functions", "delete", function).
				WithRegion(region).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithBoolFlag("quiet").
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Function deleted successfully"), nil
		},
	)

	// Call function
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_functions_call",
			Description: "Call a Cloud Function",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"function", "region"},
				"properties": map[string]any{
					"function": map[string]any{
						"type":        "string",
						"description": "Function name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region",
					},
					"data": map[string]any{
						"type":        "string",
						"description": "JSON data to send to the function",
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
			function, err := services.GetRequiredString(args, "function")
			if err != nil {
				return services.ToolError(err), nil
			}
			region, err := services.GetRequiredString(args, "region")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("functions", "call", function).
				WithRegion(region).
				WithProject(services.GetOptionalString(args, "project", ""))

			if data := services.GetOptionalString(args, "data", ""); data != "" {
				cmd.WithFlag("data", data)
			}

			result, err := cmd.WithTextFormat().ExecuteWithRegion(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.Stdout), nil
		},
	)

	// Read function logs
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_functions_logs_read",
			Description: "Read logs for a Cloud Function",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"function", "region"},
				"properties": map[string]any{
					"function": map[string]any{
						"type":        "string",
						"description": "Function name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"limit": map[string]any{
						"type":        "number",
						"description": "Maximum log entries",
						"default":     50,
					},
					"min_log_level": map[string]any{
						"type":        "string",
						"description": "Minimum log level",
						"enum":        []string{"DEBUG", "INFO", "ERROR"},
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			function, err := services.GetRequiredString(args, "function")
			if err != nil {
				return services.ToolError(err), nil
			}
			region, err := services.GetRequiredString(args, "region")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("functions", "logs", "read", function).
				WithRegion(region).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithFlag("limit", fmt.Sprintf("%d", services.GetOptionalInt(args, "limit", 50)))

			if minLevel := services.GetOptionalString(args, "min_log_level", ""); minLevel != "" {
				cmd.WithFlag("min-log-level", minLevel)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
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
