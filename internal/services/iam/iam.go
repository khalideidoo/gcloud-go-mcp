// Package iam provides MCP tools for Google Cloud IAM.
package iam

import (
	"context"
	"encoding/json"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all IAM tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List service accounts
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_iam_service_accounts_list",
			Description: "List service accounts in a project",
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
			result, err := base.Executor.Command("iam", "service-accounts", "list").
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create service account
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_iam_service_accounts_create",
			Description: "Create a service account",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"name"},
				"properties": map[string]any{
					"name": map[string]any{
						"type":        "string",
						"description": "Service account name (ID)",
					},
					"display_name": map[string]any{
						"type":        "string",
						"description": "Display name for the service account",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "Description of the service account",
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
			name, err := services.GetRequiredString(args, "name")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("iam", "service-accounts", "create", name).
				WithProject(services.GetOptionalString(args, "project", ""))

			if displayName := services.GetOptionalString(args, "display_name", ""); displayName != "" {
				cmd.WithFlag("display-name", displayName)
			}
			if description := services.GetOptionalString(args, "description", ""); description != "" {
				cmd.WithFlag("description", description)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Delete service account
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_iam_service_accounts_delete",
			Description: "Delete a service account",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"email"},
				"properties": map[string]any{
					"email": map[string]any{
						"type":        "string",
						"description": "Service account email",
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
			email, err := services.GetRequiredString(args, "email")
			if err != nil {
				return services.ToolError(err), nil
			}

			_, err = base.Executor.Command("iam", "service-accounts", "delete", email).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithBoolFlag("quiet").
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Service account deleted successfully"), nil
		},
	)

	// Describe service account
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_iam_service_accounts_describe",
			Description: "Get details of a service account",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"email"},
				"properties": map[string]any{
					"email": map[string]any{
						"type":        "string",
						"description": "Service account email",
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
			email, err := services.GetRequiredString(args, "email")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("iam", "service-accounts", "describe", email).
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// List service account keys
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_iam_service_accounts_keys_list",
			Description: "List keys for a service account",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"email"},
				"properties": map[string]any{
					"email": map[string]any{
						"type":        "string",
						"description": "Service account email",
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
			email, err := services.GetRequiredString(args, "email")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("iam", "service-accounts", "keys", "list").
				WithFlag("iam-account", email).
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create service account key
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_iam_service_accounts_keys_create",
			Description: "Create a new key for a service account (outputs to stdout)",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"email"},
				"properties": map[string]any{
					"email": map[string]any{
						"type":        "string",
						"description": "Service account email",
					},
					"key_file_type": map[string]any{
						"type":        "string",
						"description": "Key file type",
						"default":     "json",
						"enum":        []string{"json", "p12"},
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
			email, err := services.GetRequiredString(args, "email")
			if err != nil {
				return services.ToolError(err), nil
			}

			// Note: This outputs the key to /dev/stdout which may not work on all systems
			// For production use, consider writing to a file
			result, err := base.Executor.Command("iam", "service-accounts", "keys", "create", "/dev/stdout").
				WithFlag("iam-account", email).
				WithFlag("key-file-type", services.GetOptionalString(args, "key_file_type", "json")).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithTextFormat().
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.Stdout), nil
		},
	)

	// List roles
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_iam_roles_list",
			Description: "List IAM roles",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "Project for custom roles (leave empty for predefined roles)",
					},
					"show_deleted": map[string]any{
						"type":        "boolean",
						"description": "Include deleted roles",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			cmd := base.Executor.Command("iam", "roles", "list")

			if project := services.GetOptionalString(args, "project", ""); project != "" {
				cmd.WithProject(project)
			}
			if services.GetOptionalBool(args, "show_deleted", false) {
				cmd.WithBoolFlag("show-deleted")
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Describe role
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_iam_roles_describe",
			Description: "Get details of an IAM role",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"role"},
				"properties": map[string]any{
					"role": map[string]any{
						"type":        "string",
						"description": "Role ID (e.g., roles/viewer or projects/PROJECT/roles/ROLE)",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			role, err := services.GetRequiredString(args, "role")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("iam", "roles", "describe", role).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Get project IAM policy
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_get_iam_policy",
			Description: "Get IAM policy for a project",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project"},
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
			project, err := services.GetRequiredString(args, "project")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("projects", "get-iam-policy", project).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Add project IAM policy binding
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_add_iam_policy_binding",
			Description: "Add IAM policy binding to a project",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project", "member", "role"},
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"member": map[string]any{
						"type":        "string",
						"description": "Member to add (e.g., user:email@example.com)",
					},
					"role": map[string]any{
						"type":        "string",
						"description": "Role to grant (e.g., roles/viewer)",
					},
					"condition": map[string]any{
						"type":        "string",
						"description": "IAM condition expression",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			project, err := services.GetRequiredString(args, "project")
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

			cmd := base.Executor.Command("projects", "add-iam-policy-binding", project).
				WithFlag("member", member).
				WithFlag("role", role)

			if condition := services.GetOptionalString(args, "condition", ""); condition != "" {
				cmd.WithFlag("condition", condition)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Remove project IAM policy binding
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_projects_remove_iam_policy_binding",
			Description: "Remove IAM policy binding from a project",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"project", "member", "role"},
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"member": map[string]any{
						"type":        "string",
						"description": "Member to remove",
					},
					"role": map[string]any{
						"type":        "string",
						"description": "Role to revoke",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			project, err := services.GetRequiredString(args, "project")
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

			result, err := base.Executor.Command("projects", "remove-iam-policy-binding", project).
				WithFlag("member", member).
				WithFlag("role", role).
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
