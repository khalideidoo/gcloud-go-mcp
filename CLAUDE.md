# CLAUDE.md

## Project Overview

**gcloud-go-mcp** is a Model Context Protocol (MCP) server written in Go that wraps the Google Cloud CLI (`gcloud`) to provide LLM-accessible tools for managing GCP resources. It enables Claude and other MCP clients to interact with 67+ tools across 12 GCP services.

## Build & Test Commands

```bash
make build          # Build binary to ./bin/gcloud-go-mcp
make test           # Run all tests with verbose output
make test-coverage  # Run tests with coverage report (generates HTML)
make lint           # Run golangci-lint
make fmt            # Format code with go fmt
make tidy           # Tidy go.mod dependencies
make run            # Build and run the server locally
make clean          # Remove build artifacts
```

## Project Structure

```
cmd/gcloud-go-mcp/main.go    # Entry point, server init, tool registration
internal/
  config/                    # Configuration from environment variables
  executor/                  # Fluent API for building/executing gcloud commands
  services/
    base.go                  # BaseService, ToolHandler interface, helper functions
    [service]/               # 12 service packages (run, secrets, iam, storage, etc.)
```

## Architecture

**Layered design:**
1. **MCP Server** (main.go) - Creates server, registers tools, handles stdio transport
2. **Services** (internal/services/) - Tool definitions and handlers per GCP service
3. **Executor** (internal/executor/) - Fluent command builder for gcloud CLI

## Code Patterns

### Tool Naming Convention
All tools follow: `gcp_{service}_{resource}_{action}`
- `gcp_run_services_list`
- `gcp_secrets_versions_access`
- `gcp_compute_instances_create`

### Command Builder (Fluent API)
```go
executor.Command("run", "services", "list").
    WithProject(project).
    WithRegion(region).
    WithFlag("limit", "100").
    ExecuteWithRegion(ctx)
```

### Tool Handler Pattern
```go
func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := parseArgs(req)
    param := services.GetRequiredString(args, "param_name")
    result, err := base.Executor.Command("service", "resource", "action").
        WithProject(project).
        Execute(ctx)
    if err != nil {
        return services.ToolError(err), nil
    }
    return services.ToolResult(result.Output), nil
}
```

### Parameter Helpers (base.go)
- `GetRequiredString(args, key)` - Required string parameter
- `GetOptionalString(args, key, default)` - Optional string with default
- `GetOptionalInt(args, key, default)` - Optional int (JSON numbers are float64)
- `GetOptionalBool(args, key, default)` - Optional boolean
- `GetOptionalStringArray(args, key)` - Optional string slice
- `GetOptionalStringMap(args, key)` - Optional map[string]string

### Response Helpers
- `services.ToolResult(text)` - Successful result
- `services.ToolError(err)` - Error result

## Adding a New Service

1. Create `internal/services/{service}/{service}.go`
2. Implement `RegisterTools(server *mcp.Server, base *services.BaseService)`
3. Define tools with InputSchema using `jsonschema.Reflect()`
4. Register handlers that use `base.Executor` for commands
5. Add `{service}.RegisterTools(server, base)` in main.go

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `GCLOUD_PROJECT` | (empty) | Default GCP project ID |
| `GCLOUD_REGION` | (empty) | Default region |
| `GCLOUD_ZONE` | (empty) | Default zone |
| `GCLOUD_PATH` | `gcloud` | Path to gcloud binary |
| `GCLOUD_TIMEOUT` | `5m` | Command timeout |

## Testing

- Tests are co-located: `foo.go` → `foo_test.go`
- Uses table-driven tests
- Run specific package: `go test -v ./internal/services/secrets/...`
