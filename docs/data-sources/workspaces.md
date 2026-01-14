---
page_title: "anthropic_workspaces Data Source"
description: |-
  Retrieves a list of all workspaces in the Anthropic organization.
---

# anthropic_workspaces

Retrieves a list of all workspaces in the Anthropic organization.

## Example Usage

```hcl
data "anthropic_workspaces" "all" {}

output "workspace_names" {
  value = [for ws in data.anthropic_workspaces.all.workspaces : ws.name]
}
```

## Argument Reference

This data source has no required arguments.

## Attribute Reference

- `workspaces` - List of workspaces. Each workspace contains:
  - `id` - The unique identifier of the workspace.
  - `name` - The name of the workspace.
  - `display_name` - The display name of the workspace.
  - `created_at` - The timestamp when the workspace was created.
  - `archived_at` - The timestamp when the workspace was archived, if applicable.
