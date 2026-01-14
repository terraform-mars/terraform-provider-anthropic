---
page_title: "anthropic_api_key Data Source"
description: |-
  Retrieves information about an existing Anthropic API key.
---

# anthropic_api_key

Retrieves information about an existing Anthropic API key.

## Example Usage

```hcl
data "anthropic_api_key" "example" {
  id = "apikey_abc123"
}

output "api_key_status" {
  value = data.anthropic_api_key.example.status
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the API key.

## Attribute Reference

- `name` - The name of the API key.
- `workspace_id` - The ID of the workspace this API key belongs to.
- `status` - The status of the API key (`active`, `inactive`, `archived`).
- `hint` - The last 4 characters of the API key for identification.
- `created_at` - The timestamp when the API key was created.
