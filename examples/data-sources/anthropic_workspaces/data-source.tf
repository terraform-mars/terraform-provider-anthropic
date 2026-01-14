# List all workspaces in the organization
data "anthropic_workspaces" "all" {}

# Output workspace names
output "workspace_names" {
  value = [for ws in data.anthropic_workspaces.all.workspaces : ws.name]
}

# Find a specific workspace by name
locals {
  production_workspace = [
    for ws in data.anthropic_workspaces.all.workspaces :
    ws if ws.name == "production"
  ][0]
}

output "production_workspace_id" {
  value = local.production_workspace.id
}
