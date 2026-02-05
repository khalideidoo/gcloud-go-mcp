# GCloud MCP Server

A Model Context Protocol (MCP) server that wraps the Google Cloud CLI (`gcloud`) to provide LLM-accessible tools for managing GCP resources.

## Features

- **67+ tools** covering 12 GCP services
- Wraps native `gcloud` CLI commands
- JSON output for structured data
- Configurable defaults for project, region, and zone

## Supported Services

| Service | Tools | Description |
|---------|-------|-------------|
| Cloud Run | 10 | Deploy and manage containerized services |
| Secret Manager | 12 | Manage secrets and versions |
| IAM | 11 | Service accounts, roles, and policies |
| Cloud Logging | 3 | Read and write logs |
| Cloud Storage | 9 | Manage buckets and objects |
| Compute Engine | 12 | Manage VM instances and disks |
| Cloud Functions | 6 | Deploy and invoke serverless functions |
| Firestore | 6 | Manage databases and indexes |
| GKE | 6 | Manage Kubernetes clusters |
| Billing | 4 | View accounts and manage budgets |
| Pub/Sub | 8 | Manage topics and subscriptions |
| Projects | 7 | Create, list, and manage GCP projects |

## Prerequisites

1. **Google Cloud SDK** - Install from https://cloud.google.com/sdk/docs/install
2. **Authentication** - Run `gcloud auth login` to authenticate
3. **Project** - Set a default project: `gcloud config set project YOUR_PROJECT_ID`

## Installation

### From Source

```bash
git clone https://github.com/khalideidoo/gcloud-go-mcp.git
cd gcloud-go-mcp
make build
```

### Install to PATH

```bash
make install
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GCLOUD_PROJECT` | Default GCP project ID | (from gcloud config) |
| `GCLOUD_REGION` | Default region | `us-east1` |
| `GCLOUD_ZONE` | Default zone | `us-east1` |
| `GCLOUD_PATH` | Path to gcloud binary | `gcloud` |
| `GCLOUD_TIMEOUT` | Command timeout | `5m` |

### Claude Desktop Configuration

Add to your Claude Desktop configuration file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "gcloud": {
      "command": "/path/to/gcloud-go-mcp",
      "env": {
        "GCLOUD_PROJECT": "your-default-gcp-project-id",
        "GCLOUD_REGION": "your-default-gcp-region",
        "GCLOUD_ZONE": "your-default-gcp-zone",
        "GCLOUD_PATH": "your-gcloud-cli-path",
        "GCLOUD_TIMEOUT": "gcloud-cli-timeout"
      }
    }
  }
}
```

## Tool Reference

### Tool Naming Convention

All tools follow the pattern: `gcp_{service}_{resource}_{action}`

Examples:
- `gcp_run_services_list` - List Cloud Run services
- `gcp_run_services_deploy` - Deploy to Cloud Run
- `gcp_secrets_versions_access` - Access a secret version
- `gcp_compute_instances_create` - Create a VM instance

### Cloud Run Tools

| Tool | Description |
|------|-------------|
| `gcp_run_services_list` | List Cloud Run services |
| `gcp_run_services_describe` | Get service details |
| `gcp_run_services_deploy` | Deploy a container image |
| `gcp_run_services_delete` | Delete a service |
| `gcp_run_services_update_traffic` | Update traffic allocation |
| `gcp_run_services_get_iam_policy` | Get IAM policy |
| `gcp_run_services_add_iam_policy_binding` | Add IAM binding |
| `gcp_run_revisions_list` | List revisions |
| `gcp_run_jobs_list` | List jobs |
| `gcp_run_jobs_execute` | Execute a job |

### Secret Manager Tools

| Tool | Description |
|------|-------------|
| `gcp_secrets_list` | List secrets |
| `gcp_secrets_create` | Create a secret |
| `gcp_secrets_describe` | Get secret details |
| `gcp_secrets_delete` | Delete a secret |
| `gcp_secrets_versions_add` | Add a version |
| `gcp_secrets_versions_access` | Access version data |
| `gcp_secrets_versions_list` | List versions |
| `gcp_secrets_versions_disable` | Disable a version |
| `gcp_secrets_versions_enable` | Enable a version |
| `gcp_secrets_versions_destroy` | Destroy a version |
| `gcp_secrets_get_iam_policy` | Get IAM policy |
| `gcp_secrets_add_iam_policy_binding` | Add IAM binding |

### IAM Tools

| Tool | Description |
|------|-------------|
| `gcp_iam_service_accounts_list` | List service accounts |
| `gcp_iam_service_accounts_create` | Create service account |
| `gcp_iam_service_accounts_delete` | Delete service account |
| `gcp_iam_service_accounts_describe` | Get SA details |
| `gcp_iam_service_accounts_keys_list` | List SA keys |
| `gcp_iam_service_accounts_keys_create` | Create SA key |
| `gcp_iam_roles_list` | List roles |
| `gcp_iam_roles_describe` | Get role details |
| `gcp_projects_get_iam_policy` | Get project IAM policy |
| `gcp_projects_add_iam_policy_binding` | Add binding |
| `gcp_projects_remove_iam_policy_binding` | Remove binding |

### Cloud Storage Tools

| Tool | Description |
|------|-------------|
| `gcp_storage_buckets_list` | List buckets |
| `gcp_storage_buckets_describe` | Get bucket details |
| `gcp_storage_buckets_create` | Create bucket |
| `gcp_storage_buckets_delete` | Delete bucket |
| `gcp_storage_objects_list` | List objects |
| `gcp_storage_objects_cat` | Display object contents |
| `gcp_storage_objects_copy` | Copy objects |
| `gcp_storage_objects_delete` | Delete objects |
| `gcp_storage_objects_signed_url` | Generate signed URL |

### Compute Engine Tools

| Tool | Description |
|------|-------------|
| `gcp_compute_instances_list` | List instances |
| `gcp_compute_instances_describe` | Get instance details |
| `gcp_compute_instances_create` | Create instance |
| `gcp_compute_instances_delete` | Delete instance |
| `gcp_compute_instances_start` | Start instance |
| `gcp_compute_instances_stop` | Stop instance |
| `gcp_compute_instances_reset` | Reset instance |
| `gcp_compute_instances_ssh_command` | Get SSH command |
| `gcp_compute_disks_list` | List disks |
| `gcp_compute_disks_create` | Create disk |
| `gcp_compute_disks_snapshot` | Create snapshot |
| `gcp_compute_snapshots_list` | List snapshots |

### Projects Tools

| Tool | Description |
|------|-------------|
| `gcp_projects_list` | List all accessible projects |
| `gcp_projects_describe` | Get project metadata |
| `gcp_projects_create` | Create a new project |
| `gcp_projects_delete` | Delete a project |
| `gcp_projects_update` | Update project name |
| `gcp_projects_undelete` | Restore a deleted project |
| `gcp_projects_get_ancestors` | Get project hierarchy |

## Usage Examples

### List Cloud Run Services

```
"List my Cloud Run services in us-central1"
```

### Deploy to Cloud Run

```
"Deploy the image gcr.io/my-project/my-app:latest to a Cloud Run service named my-service"
```

### Access a Secret

```
"Get the value of the secret named api-key"
```

### Create a VM Instance

```
"Create a VM instance named test-vm in zone us-central1-a with machine type e2-small"
```

### Read Logs

```
"Show me the last 20 error logs from Cloud Run"
```

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Lint

```bash
make lint
```

## License

MIT License
