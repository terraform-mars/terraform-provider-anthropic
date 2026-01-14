---
page_title: "anthropic_workspace Resource"
description: |-
  Manages an Anthropic workspace.
---

# anthropic_workspace

Manages an Anthropic workspace. Workspaces allow you to organize API keys and control access to your Anthropic resources.

~> **Note:** Workspaces cannot be deleted, only archived. When this resource is destroyed, the workspace will be archived.

## Example Usage

```hcl
resource "anthropic_workspace" "example" {
  name = "production"
}
```

## Argument Reference

- `name` - (Required) The name of the workspace.

## Attribute Reference

- `id` - The unique identifier of the workspace.
- `display_name` - The display name of the workspace.
- `created_at` - The timestamp when the workspace was created.
- `archived_at` - The timestamp when the workspace was archived, if applicable.

## Import

Workspaces can be imported using the workspace ID:

```shell
terraform import anthropic_workspace.example wrkspc_abc123
```
