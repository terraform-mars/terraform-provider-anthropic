---
page_title: "anthropic_api_key Resource"
description: |-
  Manages an Anthropic API key.
---

# anthropic_api_key

Manages an Anthropic API key. API keys are used to authenticate requests to the Anthropic API and can be scoped to specific workspaces.

!> **Important:** The `key` attribute is only available immediately after creation. Store it securely as it cannot be retrieved later.

## Example Usage

### Organization-wide API Key

```hcl
resource "anthropic_api_key" "org_wide" {
  name = "org-api-key"
}
```

### Workspace-scoped API Key

```hcl
resource "anthropic_workspace" "example" {
  name = "production"
}

resource "anthropic_api_key" "workspace_key" {
  name         = "workspace-api-key"
  workspace_id = anthropic_workspace.example.id
}

output "api_key" {
  value     = anthropic_api_key.workspace_key.key
  sensitive = true
}
```

## Argument Reference

- `name` - (Required) The name of the API key.
- `workspace_id` - (Optional) The ID of the workspace this API key belongs to. If not specified, the key is organization-wide. Forces new resource if changed.
- `status` - (Optional) The status of the API key (`active`, `inactive`).

## Attribute Reference

- `id` - The unique identifier of the API key.
- `hint` - The last 4 characters of the API key for identification.
- `key` - (Sensitive) The full API key value. Only available immediately after creation.
- `created_at` - The timestamp when the API key was created.

## Import

API keys can be imported using the API key ID:

```shell
terraform import anthropic_api_key.example apikey_abc123
```

~> **Note:** When importing, the `key` attribute will not be available as it is only provided at creation time.
