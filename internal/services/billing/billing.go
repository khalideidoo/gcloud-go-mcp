// Package billing provides MCP tools for Google Cloud Billing.
package billing

import (
	"context"
	"encoding/json"
	"fmt"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Billing tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List billing accounts
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_billing_accounts_list",
			Description: "List billing accounts",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := base.Executor.Command("billing", "accounts", "list").
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Describe billing account
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_billing_accounts_describe",
			Description: "Get details of a billing account",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"account"},
				"properties": map[string]any{
					"account": map[string]any{
						"type":        "string",
						"description": "Billing account ID",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			account, err := services.GetRequiredString(args, "account")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("billing", "accounts", "describe", account).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// List budgets
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_billing_budgets_list",
			Description: "List budgets for a billing account",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"billing_account"},
				"properties": map[string]any{
					"billing_account": map[string]any{
						"type":        "string",
						"description": "Billing account ID",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			billingAccount, err := services.GetRequiredString(args, "billing_account")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("billing", "budgets", "list").
				WithFlag("billing-account", billingAccount).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create budget
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_billing_budgets_create",
			Description: "Create a budget for a billing account",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"billing_account", "display_name", "budget_amount"},
				"properties": map[string]any{
					"billing_account": map[string]any{
						"type":        "string",
						"description": "Billing account ID",
					},
					"display_name": map[string]any{
						"type":        "string",
						"description": "Display name for the budget",
					},
					"budget_amount": map[string]any{
						"type":        "string",
						"description": "Budget amount (e.g., 1000.00USD)",
					},
					"threshold_rules": map[string]any{
						"type":        "array",
						"description": "Threshold percentages for alerts (e.g., [0.5, 0.9, 1.0])",
						"items":       map[string]any{"type": "number"},
					},
					"filter_projects": map[string]any{
						"type":        "array",
						"description": "Project IDs to include in the budget",
						"items":       map[string]any{"type": "string"},
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			billingAccount, err := services.GetRequiredString(args, "billing_account")
			if err != nil {
				return services.ToolError(err), nil
			}
			displayName, err := services.GetRequiredString(args, "display_name")
			if err != nil {
				return services.ToolError(err), nil
			}
			budgetAmount, err := services.GetRequiredString(args, "budget_amount")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("billing", "budgets", "create").
				WithFlag("billing-account", billingAccount).
				WithFlag("display-name", displayName).
				WithFlag("budget-amount", budgetAmount)

			// Add threshold rules if provided
			if thresholds, ok := args["threshold_rules"].([]any); ok {
				for _, t := range thresholds {
					if threshold, ok := t.(float64); ok {
						cmd.WithArrayFlag("threshold-rule", fmt.Sprintf("percent=%g", threshold))
					}
				}
			}

			// Add project filter if provided
			if projects := services.GetOptionalStringArray(args, "filter_projects"); len(projects) > 0 {
				for _, p := range projects {
					cmd.WithArrayFlag("filter-projects", p)
				}
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
