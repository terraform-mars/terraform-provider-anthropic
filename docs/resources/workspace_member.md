---
page_title: "anthropic_workspace_member Resource"
description: |-
  Manages a member's access to an Anthropic workspace.
---

# anthropic_workspace_member

Manages a member's access to an Anthropic workspace. This resource adds users to workspaces and controls their role within that workspace.

## Example Usage

```hcl
resource "anthropic_workspace" "example" {
  name = "production"
}

resource "anthropic_workspace_member" "developer" {
  workspace_id   = anthropic_workspace.example.id
  user_id        = "user_abc123"
  workspace_role = "workspace_developer"
}
```

## Argument Reference

- `workspace_id` - (Required) The ID of the workspace. Forces new resource if changed.
- `user_id` - (Required) The ID of the user to add to the workspace. Forces new resource if changed.
- `workspace_role` - (Required) The role of the user in the workspace. Valid values:
  - `workspace_user` - Basic workspace access
  - `workspace_admin` - Administrative access to the workspace
  - `workspace_developer` - Developer access to the workspace

## Attribute Reference

- `id` - The composite identifier of the workspace member (`workspace_id/user_id`).

## Import

Workspace members can be imported using the format `workspace_id/user_id`:

```shell
terraform import anthropic_workspace_member.example wrkspc_abc123/user_xyz789
```
