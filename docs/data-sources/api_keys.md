---
page_title: "anthropic_api_keys Data Source"
description: |-
  Retrieves a list of API keys in the Anthropic organization.
---

# anthropic_api_keys

Retrieves a list of API keys in the Anthropic organization, optionally filtered by workspace or status.

## Example Usage

### List All API Keys

```hcl
data "anthropic_api_keys" "all" {}

output "api_key_names" {
  value = [for k in data.anthropic_api_keys.all.api_keys : k.name]
}
```

### Filter by Workspace

```hcl
data "anthropic_api_keys" "workspace_keys" {
  workspace_id = "wrkspc_abc123"
}
```

### Filter by Status

```hcl
data "anthropic_api_keys" "active_keys" {
  status = "active"
}
```

## Argument Reference

- `workspace_id` - (Optional) Filter API keys by workspace ID.
- `status` - (Optional) Filter API keys by status (`active`, `inactive`, `archived`).

## Attribute Reference

- `api_keys` - List of API keys. Each API key contains:
  - `id` - The unique identifier of the API key.
  - `name` - The name of the API key.
  - `workspace_id` - The ID of the workspace this API key belongs to.
  - `status` - The status of the API key.
  - `hint` - The last 4 characters of the API key.
  - `created_at` - The timestamp when the API key was created.
