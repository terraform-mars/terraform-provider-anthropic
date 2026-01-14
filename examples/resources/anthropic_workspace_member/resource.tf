# Add a user to a workspace with a specific role
resource "anthropic_workspace_member" "developer" {
  workspace_id   = anthropic_workspace.production.id
  user_id        = "user_abc123"
  workspace_role = "workspace_developer"
}

# Add an admin to a workspace
resource "anthropic_workspace_member" "admin" {
  workspace_id   = anthropic_workspace.production.id
  user_id        = "user_xyz789"
  workspace_role = "workspace_admin"
}

# Add multiple users to a workspace
locals {
  team_members = {
    "user_001" = "workspace_developer"
    "user_002" = "workspace_developer"
    "user_003" = "workspace_admin"
  }
}

resource "anthropic_workspace_member" "team" {
  for_each = local.team_members

  workspace_id   = anthropic_workspace.development.id
  user_id        = each.key
  workspace_role = each.value
}
