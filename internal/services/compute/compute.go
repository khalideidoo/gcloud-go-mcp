// Package compute provides MCP tools for Google Compute Engine.
package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Compute Engine tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List instances
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_instances_list",
			Description: "List Compute Engine VM instances",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone (leave empty for all zones)",
					},
					"filter": map[string]any{
						"type":        "string",
						"description": "Filter expression",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)

			cmd := base.Executor.Command("compute", "instances", "list").
				WithProject(services.GetOptionalString(args, "project", ""))

			if zone := services.GetOptionalString(args, "zone", ""); zone != "" {
				cmd.WithFlag("zones", zone)
			}
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

	// Describe instance
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_instances_describe",
			Description: "Get details of a VM instance",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"instance", "zone"},
				"properties": map[string]any{
					"instance": map[string]any{
						"type":        "string",
						"description": "Instance name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone of the instance",
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
			instance, err := services.GetRequiredString(args, "instance")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("compute", "instances", "describe", instance).
				WithZone(zone).
				WithProject(services.GetOptionalString(args, "project", "")).
				ExecuteWithZone(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create instance
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_instances_create",
			Description: "Create a new VM instance",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"instance", "zone"},
				"properties": map[string]any{
					"instance": map[string]any{
						"type":        "string",
						"description": "Instance name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone for the instance",
					},
					"machine_type": map[string]any{
						"type":        "string",
						"description": "Machine type (e.g., e2-micro, n1-standard-1)",
						"default":     "e2-micro",
					},
					"image_family": map[string]any{
						"type":        "string",
						"description": "Image family (e.g., debian-11, ubuntu-2204-lts)",
						"default":     "debian-11",
					},
					"image_project": map[string]any{
						"type":        "string",
						"description": "Image project",
						"default":     "debian-cloud",
					},
					"boot_disk_size": map[string]any{
						"type":        "string",
						"description": "Boot disk size (e.g., 10GB, 50GB)",
					},
					"boot_disk_type": map[string]any{
						"type":        "string",
						"description": "Boot disk type (pd-standard, pd-ssd, pd-balanced)",
					},
					"network": map[string]any{
						"type":        "string",
						"description": "Network name",
					},
					"subnet": map[string]any{
						"type":        "string",
						"description": "Subnet name",
					},
					"service_account": map[string]any{
						"type":        "string",
						"description": "Service account email",
					},
					"scopes": map[string]any{
						"type":        "array",
						"description": "API scopes",
						"items":       map[string]any{"type": "string"},
					},
					"tags": map[string]any{
						"type":        "array",
						"description": "Network tags",
						"items":       map[string]any{"type": "string"},
					},
					"labels": map[string]any{
						"type":        "object",
						"description": "Labels",
					},
					"metadata": map[string]any{
						"type":        "object",
						"description": "Metadata key-value pairs",
					},
					"preemptible": map[string]any{
						"type":        "boolean",
						"description": "Use preemptible VM",
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
			instance, err := services.GetRequiredString(args, "instance")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("compute", "instances", "create", instance).
				WithZone(zone).
				WithProject(services.GetOptionalString(args, "project", ""))

			cmd.WithFlag("machine-type", services.GetOptionalString(args, "machine_type", "e2-micro"))
			cmd.WithFlag("image-family", services.GetOptionalString(args, "image_family", "debian-11"))
			cmd.WithFlag("image-project", services.GetOptionalString(args, "image_project", "debian-cloud"))

			if bootDiskSize := services.GetOptionalString(args, "boot_disk_size", ""); bootDiskSize != "" {
				cmd.WithFlag("boot-disk-size", bootDiskSize)
			}
			if bootDiskType := services.GetOptionalString(args, "boot_disk_type", ""); bootDiskType != "" {
				cmd.WithFlag("boot-disk-type", bootDiskType)
			}
			if network := services.GetOptionalString(args, "network", ""); network != "" {
				cmd.WithFlag("network", network)
			}
			if subnet := services.GetOptionalString(args, "subnet", ""); subnet != "" {
				cmd.WithFlag("subnet", subnet)
			}
			if sa := services.GetOptionalString(args, "service_account", ""); sa != "" {
				cmd.WithFlag("service-account", sa)
			}
			if scopes := services.GetOptionalStringArray(args, "scopes"); len(scopes) > 0 {
				cmd.WithFlag("scopes", strings.Join(scopes, ","))
			}
			if tags := services.GetOptionalStringArray(args, "tags"); len(tags) > 0 {
				cmd.WithFlag("tags", strings.Join(tags, ","))
			}
			if labels := services.GetOptionalStringMap(args, "labels"); len(labels) > 0 {
				var pairs []string
				for k, v := range labels {
					pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
				}
				cmd.WithFlag("labels", strings.Join(pairs, ","))
			}
			if metadata := services.GetOptionalStringMap(args, "metadata"); len(metadata) > 0 {
				var pairs []string
				for k, v := range metadata {
					pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
				}
				cmd.WithFlag("metadata", strings.Join(pairs, ","))
			}
			if services.GetOptionalBool(args, "preemptible", false) {
				cmd.WithBoolFlag("preemptible")
			}

			result, err := cmd.ExecuteWithZone(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Delete instance
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_instances_delete",
			Description: "Delete a VM instance",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"instance", "zone"},
				"properties": map[string]any{
					"instance": map[string]any{
						"type":        "string",
						"description": "Instance name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone of the instance",
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
			instance, err := services.GetRequiredString(args, "instance")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}

			_, err = base.Executor.Command("compute", "instances", "delete", instance).
				WithZone(zone).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithBoolFlag("quiet").
				ExecuteWithZone(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Instance deleted successfully"), nil
		},
	)

	// Start instance
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_instances_start",
			Description: "Start a stopped VM instance",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"instance", "zone"},
				"properties": map[string]any{
					"instance": map[string]any{
						"type":        "string",
						"description": "Instance name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone of the instance",
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
			instance, err := services.GetRequiredString(args, "instance")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("compute", "instances", "start", instance).
				WithZone(zone).
				WithProject(services.GetOptionalString(args, "project", "")).
				ExecuteWithZone(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Stop instance
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_instances_stop",
			Description: "Stop a running VM instance",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"instance", "zone"},
				"properties": map[string]any{
					"instance": map[string]any{
						"type":        "string",
						"description": "Instance name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone of the instance",
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
			instance, err := services.GetRequiredString(args, "instance")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("compute", "instances", "stop", instance).
				WithZone(zone).
				WithProject(services.GetOptionalString(args, "project", "")).
				ExecuteWithZone(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Reset instance
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_instances_reset",
			Description: "Reset (hard reboot) a VM instance",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"instance", "zone"},
				"properties": map[string]any{
					"instance": map[string]any{
						"type":        "string",
						"description": "Instance name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone of the instance",
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
			instance, err := services.GetRequiredString(args, "instance")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("compute", "instances", "reset", instance).
				WithZone(zone).
				WithProject(services.GetOptionalString(args, "project", "")).
				ExecuteWithZone(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// SSH command
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_instances_ssh_command",
			Description: "Get SSH command for connecting to an instance",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"instance", "zone"},
				"properties": map[string]any{
					"instance": map[string]any{
						"type":        "string",
						"description": "Instance name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone of the instance",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"user": map[string]any{
						"type":        "string",
						"description": "SSH username",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			instance, err := services.GetRequiredString(args, "instance")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}
			project := services.GetOptionalString(args, "project", "")
			user := services.GetOptionalString(args, "user", "")

			// Build SSH command string
			sshCmd := fmt.Sprintf("gcloud compute ssh %s --zone=%s", instance, zone)
			if project != "" {
				sshCmd += fmt.Sprintf(" --project=%s", project)
			}
			if user != "" {
				sshCmd = fmt.Sprintf("gcloud compute ssh %s@%s --zone=%s", user, instance, zone)
				if project != "" {
					sshCmd += fmt.Sprintf(" --project=%s", project)
				}
			}

			return services.ToolResult(fmt.Sprintf("SSH command:\n%s", sshCmd)), nil
		},
	)

	// List disks
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_disks_list",
			Description: "List persistent disks",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)

			cmd := base.Executor.Command("compute", "disks", "list").
				WithProject(services.GetOptionalString(args, "project", ""))

			if zone := services.GetOptionalString(args, "zone", ""); zone != "" {
				cmd.WithFlag("zones", zone)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create disk
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_disks_create",
			Description: "Create a persistent disk",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"disk", "zone"},
				"properties": map[string]any{
					"disk": map[string]any{
						"type":        "string",
						"description": "Disk name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone",
					},
					"size": map[string]any{
						"type":        "string",
						"description": "Disk size (e.g., 100GB)",
					},
					"type": map[string]any{
						"type":        "string",
						"description": "Disk type (pd-standard, pd-ssd, pd-balanced)",
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
			disk, err := services.GetRequiredString(args, "disk")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("compute", "disks", "create", disk).
				WithZone(zone).
				WithProject(services.GetOptionalString(args, "project", ""))

			if size := services.GetOptionalString(args, "size", ""); size != "" {
				cmd.WithFlag("size", size)
			}
			if diskType := services.GetOptionalString(args, "type", ""); diskType != "" {
				cmd.WithFlag("type", diskType)
			}

			result, err := cmd.ExecuteWithZone(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// List snapshots
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_snapshots_list",
			Description: "List disk snapshots",
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

			result, err := base.Executor.Command("compute", "snapshots", "list").
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create snapshot
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_compute_disks_snapshot",
			Description: "Create a snapshot of a disk",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"disk", "zone", "snapshot_name"},
				"properties": map[string]any{
					"disk": map[string]any{
						"type":        "string",
						"description": "Disk name",
					},
					"zone": map[string]any{
						"type":        "string",
						"description": "Zone of the disk",
					},
					"snapshot_name": map[string]any{
						"type":        "string",
						"description": "Name for the snapshot",
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
			disk, err := services.GetRequiredString(args, "disk")
			if err != nil {
				return services.ToolError(err), nil
			}
			zone, err := services.GetRequiredString(args, "zone")
			if err != nil {
				return services.ToolError(err), nil
			}
			snapshotName, err := services.GetRequiredString(args, "snapshot_name")
			if err != nil {
				return services.ToolError(err), nil
			}

			result, err := base.Executor.Command("compute", "disks", "snapshot", disk).
				WithFlag("snapshot-names", snapshotName).
				WithZone(zone).
				WithProject(services.GetOptionalString(args, "project", "")).
				ExecuteWithZone(ctx)

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
