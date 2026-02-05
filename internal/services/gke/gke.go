// Package gke provides MCP tools for Google Kubernetes Engine.
package gke

import (
	"context"
	"encoding/json"
	"fmt"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all GKE tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List clusters
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_gke_clusters_list",
			Description: "List GKE clusters",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region (leave empty for all regions)",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)

			cmd := base.Executor.Command("container", "clusters", "list").
				WithProject(services.GetOptionalString(args, "project", ""))

			if region := services.GetOptionalString(args, "region", ""); region != "" {
				cmd.WithFlag("region", region)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Describe cluster
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_gke_clusters_describe",
			Description: "Get details of a GKE cluster",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"cluster"},
				"properties": map[string]any{
					"cluster": map[string]any{
						"type":        "string",
						"description": "Cluster name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region (for regional clusters)",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone (for zonal clusters)",
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
			cluster, err := services.GetRequiredString(args, "cluster")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("container", "clusters", "describe", cluster).
				WithProject(services.GetOptionalString(args, "project", ""))

			if region := services.GetOptionalString(args, "region", ""); region != "" {
				cmd.WithFlag("region", region)
			}
			if zone := services.GetOptionalString(args, "zone", ""); zone != "" {
				cmd.WithFlag("zone", zone)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create cluster
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_gke_clusters_create",
			Description: "Create a GKE cluster",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"cluster"},
				"properties": map[string]any{
					"cluster": map[string]any{
						"type":        "string",
						"description": "Cluster name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region (for regional cluster)",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone (for zonal cluster)",
					},
					"machine_type": map[string]any{
						"type":        "string",
						"description": "Machine type for nodes",
						"default":     "e2-medium",
					},
					"num_nodes": map[string]any{
						"type":        "number",
						"description": "Number of nodes",
						"default":     3,
					},
					"enable_autoscaling": map[string]any{
						"type":        "boolean",
						"description": "Enable cluster autoscaling",
					},
					"min_nodes": map[string]any{
						"type":        "number",
						"description": "Minimum nodes for autoscaling",
					},
					"max_nodes": map[string]any{
						"type":        "number",
						"description": "Maximum nodes for autoscaling",
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
			cluster, err := services.GetRequiredString(args, "cluster")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("container", "clusters", "create", cluster).
				WithProject(services.GetOptionalString(args, "project", ""))

			if region := services.GetOptionalString(args, "region", ""); region != "" {
				cmd.WithFlag("region", region)
			}
			if zone := services.GetOptionalString(args, "zone", ""); zone != "" {
				cmd.WithFlag("zone", zone)
			}

			cmd.WithFlag("machine-type", services.GetOptionalString(args, "machine_type", "e2-medium"))
			cmd.WithFlag("num-nodes", fmt.Sprintf("%d", services.GetOptionalInt(args, "num_nodes", 3)))

			if services.GetOptionalBool(args, "enable_autoscaling", false) {
				cmd.WithBoolFlag("enable-autoscaling")
				if minNodes := services.GetOptionalInt(args, "min_nodes", -1); minNodes >= 0 {
					cmd.WithFlag("min-nodes", fmt.Sprintf("%d", minNodes))
				}
				if maxNodes := services.GetOptionalInt(args, "max_nodes", -1); maxNodes >= 0 {
					cmd.WithFlag("max-nodes", fmt.Sprintf("%d", maxNodes))
				}
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Delete cluster
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_gke_clusters_delete",
			Description: "Delete a GKE cluster",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"cluster"},
				"properties": map[string]any{
					"cluster": map[string]any{
						"type":        "string",
						"description": "Cluster name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region (for regional clusters)",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone (for zonal clusters)",
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
			cluster, err := services.GetRequiredString(args, "cluster")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("container", "clusters", "delete", cluster).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithBoolFlag("quiet")

			if region := services.GetOptionalString(args, "region", ""); region != "" {
				cmd.WithFlag("region", region)
			}
			if zone := services.GetOptionalString(args, "zone", ""); zone != "" {
				cmd.WithFlag("zone", zone)
			}

			_, err = cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Cluster deleted successfully"), nil
		},
	)

	// Get credentials
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_gke_clusters_get_credentials",
			Description: "Get kubeconfig credentials for a GKE cluster",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"cluster"},
				"properties": map[string]any{
					"cluster": map[string]any{
						"type":        "string",
						"description": "Cluster name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region (for regional clusters)",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone (for zonal clusters)",
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
			cluster, err := services.GetRequiredString(args, "cluster")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("container", "clusters", "get-credentials", cluster).
				WithProject(services.GetOptionalString(args, "project", ""))

			if region := services.GetOptionalString(args, "region", ""); region != "" {
				cmd.WithFlag("region", region)
			}
			if zone := services.GetOptionalString(args, "zone", ""); zone != "" {
				cmd.WithFlag("zone", zone)
			}

			result, err := cmd.WithTextFormat().Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Credentials fetched successfully.\n" + result.Stderr), nil
		},
	)

	// List node pools
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_gke_node_pools_list",
			Description: "List node pools in a GKE cluster",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"cluster"},
				"properties": map[string]any{
					"cluster": map[string]any{
						"type":        "string",
						"description": "Cluster name",
					},
					"region": map[string]any{
						"type":        "string",
						"description": "Region (for regional clusters)",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone (for zonal clusters)",
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
			cluster, err := services.GetRequiredString(args, "cluster")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("container", "node-pools", "list").
				WithFlag("cluster", cluster).
				WithProject(services.GetOptionalString(args, "project", ""))

			if region := services.GetOptionalString(args, "region", ""); region != "" {
				cmd.WithFlag("region", region)
			}
			if zone := services.GetOptionalString(args, "zone", ""); zone != "" {
				cmd.WithFlag("zone", zone)
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
