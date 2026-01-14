# Create a workspace for each environment
resource "anthropic_workspace" "production" {
  name = "production"
}

resource "anthropic_workspace" "staging" {
  name = "staging"
}

resource "anthropic_workspace" "development" {
  name = "development"
}

# Create workspaces using for_each
resource "anthropic_workspace" "teams" {
  for_each = toset(["backend", "frontend", "data-science"])
  name     = "team-${each.key}"
}

output "production_workspace_id" {
  value = anthropic_workspace.production.id
}
