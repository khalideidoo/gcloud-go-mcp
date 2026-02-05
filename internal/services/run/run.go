// Package run provides MCP tools for Google Cloud Run.
package run

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Cloud Run tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List services
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_services_list",
			Description: "List Cloud Run services in a project",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID (uses default if not specified)",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region to list services from (uses default if not specified)",
					},
					"limit": map[string]any{
						"type":        "number",
						"description": "Maximum number of services to return",
						"default":     100,
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			project := services.GetOptionalString(args, "project", "")
			region := services.GetOptionalString(args, "region", "")
			limit := services.GetOptionalInt(args, "limit", 100)

			result, err := base.Executor.Command("run", "services", "list").
				WithProject(project).
				WithRegion(region).
				WithFlag("limit", fmt.Sprintf("%d", limit)).
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Describe service
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_services_describe",
			Description: "Get detailed information about a Cloud Run service",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"service"},
				"properties": map[string]any{
					"service": map[string]any{
						"type":        "string",
						"description": "Name of the service to describe",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region of the service",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			service, err := services.GetRequiredString(args, "service")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("run", "services", "describe", service).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", "")).
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Deploy service
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_services_deploy",
			Description: "Deploy a container image to Cloud Run",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"service", "image"},
				"properties": map[string]any{
					"service": map[string]any{
						"type":        "string",
						"description": "Name of the service to deploy",
					},
					"image": map[string]any{
						"type":        "string",
						"description": "Container image to deploy (e.g., gcr.io/project/image:tag)",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region to deploy to",
					},
					"port": map[string]any{
						"type":        "string",
						"description": "Port the container listens on",
						"default":     "8080",
					},
					"memory": map[string]any{
						"type":        "string",
						"description": "Memory limit (e.g., 512Mi, 1Gi)",
					},
					"cpu": map[string]any{
						"type":        "string",
						"description": "CPU limit (e.g., 1, 2)",
					},
					"min_instances": map[string]any{
						"type":        "number",
						"description": "Minimum number of instances",
					},
					"max_instances": map[string]any{
						"type":        "number",
						"description": "Maximum number of instances",
					},
					"service_account": map[string]any{
						"type":        "string",
						"description": "Service account email to run as",
					},
					"env_vars": map[string]any{
						"type":        "object",
						"description": "Environment variables as key-value pairs",
					},
					"allow_unauthenticated": map[string]any{
						"type":        "boolean",
						"description": "Allow unauthenticated access",
						"default":     false,
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			service, err := services.GetRequiredString(args, "service")
			if err != nil {
				return services.ToolError(err), nil
			}
			image, err := services.GetRequiredString(args, "image")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("run", "deploy", service).
				WithFlag("image", image).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", ""))

			if port := services.GetOptionalString(args, "port", ""); port != "" {
				cmd.WithFlag("port", port)
			}
			if memory := services.GetOptionalString(args, "memory", ""); memory != "" {
				cmd.WithFlag("memory", memory)
			}
			if cpu := services.GetOptionalString(args, "cpu", ""); cpu != "" {
				cmd.WithFlag("cpu", cpu)
			}
			if minInstances := services.GetOptionalInt(args, "min_instances", -1); minInstances >= 0 {
				cmd.WithFlag("min-instances", fmt.Sprintf("%d", minInstances))
			}
			if maxInstances := services.GetOptionalInt(args, "max_instances", -1); maxInstances >= 0 {
				cmd.WithFlag("max-instances", fmt.Sprintf("%d", maxInstances))
			}
			if sa := services.GetOptionalString(args, "service_account", ""); sa != "" {
				cmd.WithFlag("service-account", sa)
			}

			// Handle environment variables
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

	// Delete service
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_services_delete",
			Description: "Delete a Cloud Run service",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"service"},
				"properties": map[string]any{
					"service": map[string]any{
						"type":        "string",
						"description": "Name of the service to delete",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region of the service",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			service, err := services.GetRequiredString(args, "service")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("run", "services", "delete", service).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", "")).
				WithBoolFlag("quiet").
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Update traffic
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_services_update_traffic",
			Description: "Update traffic allocation for a Cloud Run service",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"service"},
				"properties": map[string]any{
					"service": map[string]any{
						"type":        "string",
						"description": "Name of the service",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region of the service",
					},
					"to_revisions": map[string]any{
						"type":        "string",
						"description": "Traffic allocation (e.g., 'LATEST=100' or 'rev1=50,rev2=50')",
					},
					"to_latest": map[string]any{
						"type":        "boolean",
						"description": "Send 100% traffic to latest revision",
						"default":     false,
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			service, err := services.GetRequiredString(args, "service")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("run", "services", "update-traffic", service).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", ""))

			if services.GetOptionalBool(args, "to_latest", false) {
				cmd.WithBoolFlag("to-latest")
			} else if toRevisions := services.GetOptionalString(args, "to_revisions", ""); toRevisions != "" {
				cmd.WithFlag("to-revisions", toRevisions)
			}

			result, err := cmd.ExecuteWithRegion(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Get IAM policy
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_services_get_iam_policy",
			Description: "Get IAM policy for a Cloud Run service",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"service"},
				"properties": map[string]any{
					"service": map[string]any{
						"type":        "string",
						"description": "Name of the service",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region of the service",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			service, err := services.GetRequiredString(args, "service")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("run", "services", "get-iam-policy", service).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", "")).
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Add IAM policy binding
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_services_add_iam_policy_binding",
			Description: "Add IAM policy binding to a Cloud Run service",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"service", "member", "role"},
				"properties": map[string]any{
					"service": map[string]any{
						"type":        "string",
						"description": "Name of the service",
					},
					"member": map[string]any{
						"type":        "string",
						"description": "Member to add (e.g., user:email@example.com, serviceAccount:sa@project.iam.gserviceaccount.com)",
					},
					"role": map[string]any{
						"type":        "string",
						"description": "Role to grant (e.g., roles/run.invoker)",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region of the service",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			service, err := services.GetRequiredString(args, "service")
			if err != nil {
				return services.ToolError(err), nil
			}
			member, err := services.GetRequiredString(args, "member")
			if err != nil {
				return services.ToolError(err), nil
			}
			role, err := services.GetRequiredString(args, "role")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("run", "services", "add-iam-policy-binding", service).
				WithFlag("member", member).
				WithFlag("role", role).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", "")).
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// List revisions
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_revisions_list",
			Description: "List revisions for a Cloud Run service",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"service"},
				"properties": map[string]any{
					"service": map[string]any{
						"type":        "string",
						"description": "Name of the service",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region of the service",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			service, err := services.GetRequiredString(args, "service")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("run", "revisions", "list").
				WithFlag("service", service).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", "")).
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// List jobs
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_jobs_list",
			Description: "List Cloud Run jobs",
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

			result, err := base.Executor.Command("run", "jobs", "list").
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", "")).
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Execute job
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_run_jobs_execute",
			Description: "Execute a Cloud Run job",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"job"},
				"properties": map[string]any{
					"job": map[string]any{
						"type":        "string",
						"description": "Name of the job",
					},
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
			job, err := services.GetRequiredString(args, "job")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("run", "jobs", "execute", job).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithRegion(services.GetOptionalString(args, "region", "")).
				ExecuteWithRegion(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)
}

// parseArgs extracts arguments from the request.
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
