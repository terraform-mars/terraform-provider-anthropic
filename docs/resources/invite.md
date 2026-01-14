---
page_title: "anthropic_invite Resource"
description: |-
  Manages an invitation to join the Anthropic organization.
---

# anthropic_invite

Manages an invitation to join the Anthropic organization. Invites allow you to add new users to your organization.

~> **Note:** Invites cannot be updated. Any changes require recreating the invite.

## Example Usage

```hcl
resource "anthropic_invite" "new_developer" {
  email = "developer@example.com"
  role  = "developer"
}
```

## Argument Reference

- `email` - (Required) The email address to send the invitation to. Forces new resource if changed.
- `role` - (Required) The role to assign to the invited user. Forces new resource if changed. Valid values:
  - `user` - Basic organization access
  - `admin` - Administrative access
  - `developer` - Developer access

## Attribute Reference

- `id` - The unique identifier of the invite.
- `status` - The status of the invite (`pending`, `accepted`, `expired`, `deleted`).
- `created_at` - The timestamp when the invite was created.
- `expires_at` - The timestamp when the invite expires.
- `inviter_id` - The ID of the user who created the invite.

## Import

Invites can be imported using the invite ID:

```shell
terraform import anthropic_invite.example invite_abc123
```
