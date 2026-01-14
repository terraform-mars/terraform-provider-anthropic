# Terraform Provider for Anthropic

A Terraform provider for managing Anthropic organization resources using the [Admin API](https://docs.anthropic.com/en/api/admin-api).

## Features

This provider allows you to manage:

- **Workspaces** - Organize your API keys and team access
- **API Keys** - Create and manage API keys scoped to workspaces
- **Workspace Members** - Control user access to workspaces
- **Invites** - Invite new users to your organization

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (to build the provider)
- An Anthropic Admin API key (`sk-ant-admin-...`)

## Installation

### From Terraform Registry (Recommended)

```hcl
terraform {
  required_providers {
    anthropic = {
      source  = "terraform-mars/anthropic"
      version = "~> 0.1"
    }
  }
}
```

### Building from Source

```bash
git clone https://github.com/terraform-mars/terraform-provider-anthropic.git
cd terraform-provider-anthropic
go build -o terraform-provider-anthropic
```

## Configuration

```hcl
provider "anthropic" {
  admin_key = var.anthropic_admin_key  # Or use ANTHROPIC_ADMIN_KEY env var
}
```

### Authentication

The provider requires an Admin API key, which can be provided via:

1. Provider configuration: `admin_key = "sk-ant-admin-..."`
2. Environment variable: `ANTHROPIC_ADMIN_KEY`

## Usage Examples

### Create Workspaces for Different Environments

```hcl
resource "anthropic_workspace" "envs" {
  for_each = toset(["development", "staging", "production"])
  name     = each.key
}
```

### Generate API Keys per Workspace

```hcl
resource "anthropic_api_key" "backend" {
  for_each     = anthropic_workspace.envs
  workspace_id = each.value.id
  name         = "backend-${each.key}"
}

# Output keys for CI/CD
output "api_keys" {
  sensitive = true
  value = {
    for env, key in anthropic_api_key.backend :
    env => key.key
  }
}
```

### Manage Workspace Members

```hcl
resource "anthropic_workspace_member" "developer" {
  workspace_id   = anthropic_workspace.envs["development"].id
  user_id        = "user_abc123"
  workspace_role = "workspace_developer"
}
```

### Invite New Users

```hcl
resource "anthropic_invite" "new_hire" {
  email = "newdev@company.com"
  role  = "developer"
}
```

## Resources

| Resource | Description |
|----------|-------------|
| `anthropic_workspace` | Manage workspaces |
| `anthropic_api_key` | Manage API keys |
| `anthropic_workspace_member` | Manage workspace membership |
| `anthropic_invite` | Manage organization invites |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| `anthropic_workspace` | Read a single workspace |
| `anthropic_workspaces` | List all workspaces |
| `anthropic_api_key` | Read a single API key |
| `anthropic_api_keys` | List API keys (with optional filters) |

## Development

### Building

```bash
go build -o terraform-provider-anthropic
```

### Testing

```bash
go test ./...
```

### Using Local Provider

Create a `~/.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "terraform-mars/anthropic" = "/path/to/terraform-provider-anthropic"
  }
  direct {}
}
```

## License

MIT License - see [LICENSE](LICENSE) for details.
