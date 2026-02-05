// Package firestore provides MCP tools for Google Cloud Firestore.
package firestore

import (
	"context"
	"encoding/json"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Firestore tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List databases
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_firestore_databases_list",
			Description: "List Firestore databases",
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

			result, err := base.Executor.Command("firestore", "databases", "list").
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create database
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_firestore_databases_create",
			Description: "Create a new Firestore database",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"database", "location"},
				"properties": map[string]any{
					"database": map[string]any{
						"type":        "string",
						"description": "Database ID",
					},
					"location": map[string]any{
						"type":        "string",
						"description": "Location (e.g., nam5, eur3, us-central1)",
					},
					"type": map[string]any{
						"type":        "string",
						"description": "Database type (firestore-native or datastore-mode)",
						"default":     "firestore-native",
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
			database, err := services.GetRequiredString(args, "database")
			if err != nil {
				return services.ToolError(err), nil
			}
			location, err := services.GetRequiredString(args, "location")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("firestore", "databases", "create").
				WithFlag("database", database).
				WithFlag("location", location).
				WithFlag("type", services.GetOptionalString(args, "type", "firestore-native")).
				WithProject(services.GetOptionalString(args, "project", ""))

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Describe database
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_firestore_databases_describe",
			Description: "Get details of a Firestore database",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"database"},
				"properties": map[string]any{
					"database": map[string]any{
						"type":        "string",
						"description": "Database ID",
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
			database, err := services.GetRequiredString(args, "database")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("firestore", "databases", "describe").
				WithFlag("database", database).
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Export database
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_firestore_export",
			Description: "Export Firestore data to Cloud Storage",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"output_uri_prefix"},
				"properties": map[string]any{
					"output_uri_prefix": map[string]any{
						"type":        "string",
						"description": "Cloud Storage URI prefix (gs://bucket/path)",
					},
					"database": map[string]any{
						"type":        "string",
						"description": "Database ID (default: (default))",
						"default":     "(default)",
					},
					"collection_ids": map[string]any{
						"type":        "array",
						"description": "Collection IDs to export (empty = all)",
						"items":       map[string]any{"type": "string"},
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
			outputURI, err := services.GetRequiredString(args, "output_uri_prefix")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("firestore", "export", outputURI).
				WithFlag("database", services.GetOptionalString(args, "database", "(default)")).
				WithProject(services.GetOptionalString(args, "project", ""))

			if collectionIDs := services.GetOptionalStringArray(args, "collection_ids"); len(collectionIDs) > 0 {
				for _, id := range collectionIDs {
					cmd.WithArrayFlag("collection-ids", id)
				}
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Import database
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_firestore_import",
			Description: "Import Firestore data from Cloud Storage",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"input_uri_prefix"},
				"properties": map[string]any{
					"input_uri_prefix": map[string]any{
						"type":        "string",
						"description": "Cloud Storage URI prefix (gs://bucket/path)",
					},
					"database": map[string]any{
						"type":        "string",
						"description": "Database ID (default: (default))",
						"default":     "(default)",
					},
					"collection_ids": map[string]any{
						"type":        "array",
						"description": "Collection IDs to import (empty = all)",
						"items":       map[string]any{"type": "string"},
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
			inputURI, err := services.GetRequiredString(args, "input_uri_prefix")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("firestore", "import", inputURI).
				WithFlag("database", services.GetOptionalString(args, "database", "(default)")).
				WithProject(services.GetOptionalString(args, "project", ""))

			if collectionIDs := services.GetOptionalStringArray(args, "collection_ids"); len(collectionIDs) > 0 {
				for _, id := range collectionIDs {
					cmd.WithArrayFlag("collection-ids", id)
				}
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// List indexes
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_firestore_indexes_list",
			Description: "List Firestore indexes",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"database": map[string]any{
						"type":        "string",
						"description": "Database ID (default: (default))",
						"default":     "(default)",
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

			result, err := base.Executor.Command("firestore", "indexes", "composite", "list").
				WithFlag("database", services.GetOptionalString(args, "database", "(default)")).
				WithProject(services.GetOptionalString(args, "project", "")).
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
