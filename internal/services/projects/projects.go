// Package projects provides MCP tools for Google Cloud Projects.
package projects

import (
	"context"
	"encoding/json"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Projects tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List projects
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_list",
			Description: "List all GCP projects accessible by the active account",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"filter": map[string]any{
						"type":        "string",
						"description": "Filter expression (e.g., 'name:my-project*' or 'lifecycleState:ACTIVE')",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			cmd := base.Executor.Command("projects", "list")

			if filter := services.GetOptionalString(args, "filter", ""); filter != "" {
				cmd.WithFlag("filter", filter)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Describe project
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_describe",
			Description: "Get metadata for a project",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project_id"},
				"properties": map[string]any{
					"project_id": map[string]any{
						"type":        "string",
						"description": "Project ID",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			projectID, err := services.GetRequiredString(args, "project_id")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("projects", "describe", projectID).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create project
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_create",
			Description: "Create a new GCP project",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project_id"},
				"properties": map[string]any{
					"project_id": map[string]any{
						"type":        "string",
						"description": "Unique project ID (6-30 characters, lowercase letters, digits, hyphens)",
					},
					"name": map[string]any{
						"type":        "string",
						"description": "Display name for the project",
					},
					"organization": map[string]any{
						"type":        "string",
						"description": "Organization ID to create the project under",
					},
					"folder": map[string]any{
						"type":        "string",
						"description": "Folder ID to create the project under",
					},
					"labels": map[string]any{
						"type":        "object",
						"description": "Labels to attach to the project (key-value pairs)",
						"additionalProperties": map[string]any{
							"type": "string",
						},
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			projectID, err := services.GetRequiredString(args, "project_id")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("projects", "create", projectID)

			if name := services.GetOptionalString(args, "name", ""); name != "" {
				cmd.WithFlag("name", name)
			}
			if org := services.GetOptionalString(args, "organization", ""); org != "" {
				cmd.WithFlag("organization", org)
			}
			if folder := services.GetOptionalString(args, "folder", ""); folder != "" {
				cmd.WithFlag("folder", folder)
			}
			if labels, ok := args["labels"].(map[string]any); ok && len(labels) > 0 {
				labelStr := ""
				for k, v := range labels {
					if labelStr != "" {
						labelStr += ","
					}
					labelStr += k + "=" + v.(string)
				}
				cmd.WithFlag("labels", labelStr)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Delete project
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_delete",
			Description: "Delete a project (moves to DELETE_REQUESTED state, can be restored within 30 days)",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project_id"},
				"properties": map[string]any{
					"project_id": map[string]any{
						"type":        "string",
						"description": "Project ID to delete",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			projectID, err := services.GetRequiredString(args, "project_id")
			if err != nil {
				return services.ToolError(err), nil
			}

			_, err = base.Executor.Command("projects", "delete", projectID).
				WithBoolFlag("quiet").
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Project " + projectID + " marked for deletion"), nil
		},
	)

	// Update project
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_update",
			Description: "Update the name of a project",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project_id", "name"},
				"properties": map[string]any{
					"project_id": map[string]any{
						"type":        "string",
						"description": "Project ID to update",
					},
					"name": map[string]any{
						"type":        "string",
						"description": "New display name for the project",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			projectID, err := services.GetRequiredString(args, "project_id")
			if err != nil {
				return services.ToolError(err), nil
			}
			name, err := services.GetRequiredString(args, "name")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("projects", "update", projectID).
				WithFlag("name", name).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Undelete project
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_undelete",
			Description: "Restore a project that was marked for deletion",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project_id"},
				"properties": map[string]any{
					"project_id": map[string]any{
						"type":        "string",
						"description": "Project ID to restore",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			projectID, err := services.GetRequiredString(args, "project_id")
			if err != nil {
				return services.ToolError(err), nil
			}

			_, err = base.Executor.Command("projects", "undelete", projectID).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Project " + projectID + " restored successfully"), nil
		},
	)

	// Get ancestors
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_get_ancestors",
			Description: "Get the ancestors (folder and organization hierarchy) for a project",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project_id"},
				"properties": map[string]any{
					"project_id": map[string]any{
						"type":        "string",
						"description": "Project ID",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			projectID, err := services.GetRequiredString(args, "project_id")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("projects", "get-ancestors", projectID).
				Execute(ctx)

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
