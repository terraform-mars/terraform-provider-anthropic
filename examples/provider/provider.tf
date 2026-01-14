terraform {
  required_providers {
    anthropic = {
      source  = "terraform-mars/anthropic"
      version = "~> 0.1"
    }
  }
}

# Configure the Anthropic provider
# The admin_key can also be set via the ANTHROPIC_ADMIN_KEY environment variable
provider "anthropic" {
  admin_key = var.anthropic_admin_key
}

variable "anthropic_admin_key" {
  description = "Anthropic Admin API key (sk-ant-admin-...)"
  type        = string
  sensitive   = true
}
