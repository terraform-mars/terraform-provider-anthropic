---
page_title: "Provider: Anthropic"
description: |-
  Terraform provider for managing Anthropic organization resources.
---

# Anthropic Provider

The Anthropic provider allows you to manage Anthropic organization resources such as workspaces, API keys, members, and invites using the [Admin API](https://docs.anthropic.com/en/docs/administration-api).

## Example Usage

```hcl
terraform {
  required_providers {
    anthropic = {
      source  = "terraform-mars/anthropic"
      version = "~> 0.1"
    }
  }
}

provider "anthropic" {
  admin_key = var.anthropic_admin_key
}

resource "anthropic_workspace" "example" {
  name = "my-workspace"
}

resource "anthropic_api_key" "example" {
  name         = "terraform-api-key"
  workspace_id = anthropic_workspace.example.id
}
```

## Authentication

The provider requires an Anthropic Admin API key. You can obtain one from the [Anthropic Console](https://console.anthropic.com/settings/admin-keys).

Configure the key via:
- Provider configuration: `admin_key`
- Environment variable: `ANTHROPIC_ADMIN_KEY`

## Argument Reference

- `admin_key` - (Optional) Anthropic Admin API key. Can also be set via `ANTHROPIC_ADMIN_KEY` environment variable.
- `base_url` - (Optional) Anthropic API base URL. Defaults to `https://api.anthropic.com`. Can also be set via `ANTHROPIC_BASE_URL` environment variable.
