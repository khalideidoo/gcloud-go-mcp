// Package storage provides MCP tools for Google Cloud Storage.
package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Cloud Storage tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List buckets
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_buckets_list",
			Description: "List Cloud Storage buckets",
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

			result, err := base.Executor.Command("storage", "buckets", "list").
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Describe bucket
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_buckets_describe",
			Description: "Get details of a bucket",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"bucket"},
				"properties": map[string]any{
					"bucket": map[string]any{
						"type":        "string",
						"description": "Bucket name (without gs://)",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			bucket, err := services.GetRequiredString(args, "bucket")
			if err != nil {
				return services.ToolError(err), nil
			}

			bucketURL := fmt.Sprintf("gs://%s", bucket)
			result, err := base.Executor.Command("storage", "buckets", "describe", bucketURL).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create bucket
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_buckets_create",
			Description: "Create a new bucket",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"bucket"},
				"properties": map[string]any{
					"bucket": map[string]any{
						"type":        "string",
						"description": "Bucket name",
					},
					"location": map[string]any{
						"type":        "string",
						"description": "Bucket location (e.g., US, EU, us-central1)",
						"default":     "US",
					},
					"storage_class": map[string]any{
						"type":        "string",
						"description": "Storage class (STANDARD, NEARLINE, COLDLINE, ARCHIVE)",
						"default":     "STANDARD",
					},
					"uniform_bucket_level_access": map[string]any{
						"type":        "boolean",
						"description": "Enable uniform bucket-level access",
						"default":     true,
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
			bucket, err := services.GetRequiredString(args, "bucket")
			if err != nil {
				return services.ToolError(err), nil
			}

			bucketURL := fmt.Sprintf("gs://%s", bucket)
			cmd := base.Executor.Command("storage", "buckets", "create", bucketURL).
				WithProject(services.GetOptionalString(args, "project", ""))

			if location := services.GetOptionalString(args, "location", "US"); location != "" {
				cmd.WithFlag("location", location)
			}
			if storageClass := services.GetOptionalString(args, "storage_class", "STANDARD"); storageClass != "" {
				cmd.WithFlag("default-storage-class", storageClass)
			}
			if services.GetOptionalBool(args, "uniform_bucket_level_access", true) {
				cmd.WithBoolFlag("uniform-bucket-level-access")
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Delete bucket
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_buckets_delete",
			Description: "Delete a bucket (must be empty)",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"bucket"},
				"properties": map[string]any{
					"bucket": map[string]any{
						"type":        "string",
						"description": "Bucket name",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			bucket, err := services.GetRequiredString(args, "bucket")
			if err != nil {
				return services.ToolError(err), nil
			}

			bucketURL := fmt.Sprintf("gs://%s", bucket)
			_, err = base.Executor.Command("storage", "buckets", "delete", bucketURL).
				WithTextFormat().
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Bucket deleted successfully"), nil
		},
	)

	// List objects
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_objects_list",
			Description: "List objects in a bucket",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"bucket"},
				"properties": map[string]any{
					"bucket": map[string]any{
						"type":        "string",
						"description": "Bucket name",
					},
					"prefix": map[string]any{
						"type":        "string",
						"description": "Object prefix (folder path)",
					},
					"limit": map[string]any{
						"type":        "number",
						"description": "Maximum objects to list",
						"default":     100,
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			bucket, err := services.GetRequiredString(args, "bucket")
			if err != nil {
				return services.ToolError(err), nil
			}

			bucketURL := fmt.Sprintf("gs://%s", bucket)
			if prefix := services.GetOptionalString(args, "prefix", ""); prefix != "" {
				bucketURL = fmt.Sprintf("gs://%s/%s", bucket, prefix)
			}

			// Note: gcloud storage ls doesn't have a direct limit flag
			// We use the command as-is
			result, err := base.Executor.Command("storage", "ls", bucketURL).
				WithBoolFlag("long").
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Cat object
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_objects_cat",
			Description: "Display contents of an object",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"url"},
				"properties": map[string]any{
					"url": map[string]any{
						"type":        "string",
						"description": "Object URL (gs://bucket/path/to/object)",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			url, err := services.GetRequiredString(args, "url")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("storage", "cat", url).
				WithTextFormat().
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.Stdout), nil
		},
	)

	// Copy objects
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_objects_copy",
			Description: "Copy objects between buckets or within a bucket",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"source", "destination"},
				"properties": map[string]any{
					"source": map[string]any{
						"type":        "string",
						"description": "Source URL (gs://bucket/object or local path)",
					},
					"destination": map[string]any{
						"type":        "string",
						"description": "Destination URL (gs://bucket/object or local path)",
					},
					"recursive": map[string]any{
						"type":        "boolean",
						"description": "Copy recursively",
						"default":     false,
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			source, err := services.GetRequiredString(args, "source")
			if err != nil {
				return services.ToolError(err), nil
			}
			destination, err := services.GetRequiredString(args, "destination")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("storage", "cp", source, destination)

			if services.GetOptionalBool(args, "recursive", false) {
				cmd.WithBoolFlag("recursive")
			}

			result, err := cmd.WithTextFormat().Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			if result.Stdout == "" {
				return services.ToolResult("Copy completed successfully"), nil
			}
			return services.ToolResult(result.Stdout), nil
		},
	)

	// Delete objects
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_objects_delete",
			Description: "Delete objects",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"url"},
				"properties": map[string]any{
					"url": map[string]any{
						"type":        "string",
						"description": "Object URL (gs://bucket/path/to/object)",
					},
					"recursive": map[string]any{
						"type":        "boolean",
						"description": "Delete recursively",
						"default":     false,
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			url, err := services.GetRequiredString(args, "url")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("storage", "rm", url)

			if services.GetOptionalBool(args, "recursive", false) {
				cmd.WithBoolFlag("recursive")
			}

			_, err = cmd.WithTextFormat().Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Delete completed successfully"), nil
		},
	)

	// Generate signed URL
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_storage_objects_signed_url",
			Description: "Generate a signed URL for an object",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"url"},
				"properties": map[string]any{
					"url": map[string]any{
						"type":        "string",
						"description": "Object URL (gs://bucket/path/to/object)",
					},
					"duration": map[string]any{
						"type":        "string",
						"description": "URL validity duration (e.g., 1h, 30m)",
						"default":     "1h",
					},
					"http_method": map[string]any{
						"type":        "string",
						"description": "HTTP method for the signed URL",
						"default":     "GET",
						"enum":        []string{"GET", "PUT", "DELETE"},
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			url, err := services.GetRequiredString(args, "url")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("storage", "sign-url", url).
				WithFlag("duration", services.GetOptionalString(args, "duration", "1h")).
				WithFlag("http-verb", services.GetOptionalString(args, "http_method", "GET")).
				WithTextFormat().
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
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
