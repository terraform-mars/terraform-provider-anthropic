---
page_title: "anthropic_workspace Data Source"
description: |-
  Retrieves information about an existing Anthropic workspace.
---

# anthropic_workspace

Retrieves information about an existing Anthropic workspace.

## Example Usage

```hcl
data "anthropic_workspace" "example" {
  id = "wrkspc_abc123"
}

output "workspace_name" {
  value = data.anthropic_workspace.example.name
}
```

## Argument Reference

- `id` - (Required) The unique identifier of the workspace.

## Attribute Reference

- `name` - The name of the workspace.
- `display_name` - The display name of the workspace.
- `created_at` - The timestamp when the workspace was created.
- `archived_at` - The timestamp when the workspace was archived, if applicable.
