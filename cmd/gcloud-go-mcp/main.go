// Package main is the entry point for the gcloud MCP server.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gcloud-go-mcp/internal/config"
	"gcloud-go-mcp/internal/services"
	"gcloud-go-mcp/internal/services/billing"
	"gcloud-go-mcp/internal/services/compute"
	"gcloud-go-mcp/internal/services/firestore"
	"gcloud-go-mcp/internal/services/functions"
	"gcloud-go-mcp/internal/services/gke"
	"gcloud-go-mcp/internal/services/iam"
	"gcloud-go-mcp/internal/services/logging"
	"gcloud-go-mcp/internal/services/projects"
	"gcloud-go-mcp/internal/services/pubsub"
	"gcloud-go-mcp/internal/services/run"
	"gcloud-go-mcp/internal/services/secrets"
	"gcloud-go-mcp/internal/services/storage"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "gcloud-go-mcp"
	serverVersion = "1.0.0"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create MCP server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    serverName,
			Version: serverVersion,
		},
		&mcp.ServerOptions{
			Instructions: `GCloud MCP Server provides tools for managing Google Cloud Platform resources.

Before using tools, ensure:
1. gcloud CLI is installed and configured
2. You are authenticated (gcloud auth login)
3. A default project is set, or specify project in each tool call

Tools follow the pattern: gcp_{service}_{resource}_{action}

Example usage:
- List Cloud Run services: gcp_run_services_list
- Deploy to Cloud Run: gcp_run_services_deploy
- Read logs: gcp_logging_read
- Manage secrets: gcp_secrets_create, gcp_secrets_versions_access`,
		},
	)

	// Create base service with shared executor
	base := services.NewBaseService(cfg)

	// Register all service tools
	run.RegisterTools(server, base)
	secrets.RegisterTools(server, base)
	iam.RegisterTools(server, base)
	logging.RegisterTools(server, base)
	storage.RegisterTools(server, base)
	compute.RegisterTools(server, base)
	functions.RegisterTools(server, base)
	firestore.RegisterTools(server, base)
	gke.RegisterTools(server, base)
	billing.RegisterTools(server, base)
	pubsub.RegisterTools(server, base)
	projects.RegisterTools(server, base)

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	// Run server on stdio transport
	log.Printf("Starting %s v%s", serverName, serverVersion)
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
