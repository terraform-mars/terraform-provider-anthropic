# Create an API key for a specific workspace
resource "anthropic_api_key" "backend_prod" {
  name         = "backend-production"
  workspace_id = anthropic_workspace.production.id
}

# Create an organization-wide API key (no workspace)
resource "anthropic_api_key" "admin" {
  name = "admin-key"
}

# Create API keys for each environment
resource "anthropic_api_key" "env_keys" {
  for_each = {
    dev     = anthropic_workspace.development.id
    staging = anthropic_workspace.staging.id
    prod    = anthropic_workspace.production.id
  }

  name         = "backend-${each.key}"
  workspace_id = each.value
}

# Output API keys for use in CI/CD (marked sensitive)
output "api_keys" {
  description = "API keys for each environment"
  sensitive   = true
  value = {
    for env, key in anthropic_api_key.env_keys :
    env => key.key
  }
}
