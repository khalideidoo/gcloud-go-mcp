// Package pubsub provides MCP tools for Google Cloud Pub/Sub.
package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"gcloud-go-mcp/internal/services"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterTools registers all Pub/Sub tools with the MCP server.
func RegisterTools(server *mcp.Server, base *services.BaseService) {
	// List topics
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_pubsub_topics_list",
			Description: "List Pub/Sub topics",
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

			result, err := base.Executor.Command("pubsub", "topics", "list").
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create topic
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_pubsub_topics_create",
			Description: "Create a Pub/Sub topic",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"topic"},
				"properties": map[string]any{
					"topic": map[string]any{
						"type":        "string",
						"description": "Topic name",
					},
					"project": map[string]any{
						"type":        "string",
						"description": "GCP project ID",
					},
					"labels": map[string]any{
						"type":        "object",
						"description": "Labels for the topic",
					},
				},
			},
		},
		func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := parseArgs(req)
			topic, err := services.GetRequiredString(args, "topic")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("pubsub", "topics", "create", topic).
				WithProject(services.GetOptionalString(args, "project", ""))

			if labels := services.GetOptionalStringMap(args, "labels"); len(labels) > 0 {
				var pairs []string
				for k, v := range labels {
					pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
				}
				cmd.WithFlag("labels", fmt.Sprintf("%v", pairs))
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Delete topic
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_pubsub_topics_delete",
			Description: "Delete a Pub/Sub topic",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"topic"},
				"properties": map[string]any{
					"topic": map[string]any{
						"type":        "string",
						"description": "Topic name",
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
			topic, err := services.GetRequiredString(args, "topic")
			if err != nil {
				return services.ToolError(err), nil
			}

			_, err = base.Executor.Command("pubsub", "topics", "delete", topic).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithBoolFlag("quiet").
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Topic deleted successfully"), nil
		},
	)

	// Publish message
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_pubsub_topics_publish",
			Description: "Publish a message to a Pub/Sub topic",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"topic", "message"},
				"properties": map[string]any{
					"topic": map[string]any{
						"type":        "string",
						"description": "Topic name",
					},
					"message": map[string]any{
						"type":        "string",
						"description": "Message to publish",
					},
					"attributes": map[string]any{
						"type":        "object",
						"description": "Message attributes as key-value pairs",
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
			topic, err := services.GetRequiredString(args, "topic")
			if err != nil {
				return services.ToolError(err), nil
			}
			message, err := services.GetRequiredString(args, "message")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("pubsub", "topics", "publish", topic).
				WithFlag("message", message).
				WithProject(services.GetOptionalString(args, "project", ""))

			if attrs := services.GetOptionalStringMap(args, "attributes"); len(attrs) > 0 {
				for k, v := range attrs {
					cmd.WithArrayFlag("attribute", fmt.Sprintf("%s=%s", k, v))
				}
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// List subscriptions
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_pubsub_subscriptions_list",
			Description: "List Pub/Sub subscriptions",
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

			result, err := base.Executor.Command("pubsub", "subscriptions", "list").
				WithProject(services.GetOptionalString(args, "project", "")).
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Create subscription
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_pubsub_subscriptions_create",
			Description: "Create a Pub/Sub subscription",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"subscription", "topic"},
				"properties": map[string]any{
					"subscription": map[string]any{
						"type":        "string",
						"description": "Subscription name",
					},
					"topic": map[string]any{
						"type":        "string",
						"description": "Topic to subscribe to",
					},
					"ack_deadline": map[string]any{
						"type":        "number",
						"description": "Acknowledgement deadline in seconds",
						"default":     10,
					},
					"push_endpoint": map[string]any{
						"type":        "string",
						"description": "Push endpoint URL (for push subscriptions)",
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
			subscription, err := services.GetRequiredString(args, "subscription")
			if err != nil {
				return services.ToolError(err), nil
			}
			topic, err := services.GetRequiredString(args, "topic")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("pubsub", "subscriptions", "create", subscription).
				WithFlag("topic", topic).
				WithProject(services.GetOptionalString(args, "project", ""))

			if ackDeadline := services.GetOptionalInt(args, "ack_deadline", 10); ackDeadline > 0 {
				cmd.WithFlag("ack-deadline", fmt.Sprintf("%d", ackDeadline))
			}
			if pushEndpoint := services.GetOptionalString(args, "push_endpoint", ""); pushEndpoint != "" {
				cmd.WithFlag("push-endpoint", pushEndpoint)
			}

			result, err := cmd.Execute(ctx)
			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult(result.ToJSONString()), nil
		},
	)

	// Delete subscription
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_pubsub_subscriptions_delete",
			Description: "Delete a Pub/Sub subscription",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"subscription"},
				"properties": map[string]any{
					"subscription": map[string]any{
						"type":        "string",
						"description": "Subscription name",
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
			subscription, err := services.GetRequiredString(args, "subscription")
			if err != nil {
				return services.ToolError(err), nil
			}

			_, err = base.Executor.Command("pubsub", "subscriptions", "delete", subscription).
				WithProject(services.GetOptionalString(args, "project", "")).
				WithBoolFlag("quiet").
				Execute(ctx)

			if err != nil {
				return services.ToolError(err), nil
			}
			return services.ToolResult("Subscription deleted successfully"), nil
		},
	)

	// Pull messages
	server.AddTool(
		&mcp.Tool{
			Name:        "gcp_pubsub_subscriptions_pull",
			Description: "Pull messages from a Pub/Sub subscription",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"subscription"},
				"properties": map[string]any{
					"subscription": map[string]any{
						"type":        "string",
						"description": "Subscription name",
					},
					"limit": map[string]any{
						"type":        "number",
						"description": "Maximum messages to pull",
						"default":     10,
					},
					"auto_ack": map[string]any{
						"type":        "boolean",
						"description": "Automatically acknowledge messages",
						"default":     false,
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
			subscription, err := services.GetRequiredString(args, "subscription")
			if err != nil {
				return services.ToolError(err), nil
			}

			cmd := base.Executor.Command("pubsub", "subscriptions", "pull", subscription).
				WithFlag("limit", fmt.Sprintf("%d", services.GetOptionalInt(args, "limit", 10))).
				WithProject(services.GetOptionalString(args, "project", ""))

			if services.GetOptionalBool(args, "auto_ack", false) {
				cmd.WithBoolFlag("auto-ack")
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
