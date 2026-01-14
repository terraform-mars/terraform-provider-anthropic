# Invite a new user to the organization
resource "anthropic_invite" "new_developer" {
  email = "developer@example.com"
  role  = "developer"
}

# Invite an admin
resource "anthropic_invite" "new_admin" {
  email = "admin@example.com"
  role  = "admin"
}

# Invite multiple users
resource "anthropic_invite" "batch" {
  for_each = {
    "alice@example.com" = "developer"
    "bob@example.com"   = "developer"
    "carol@example.com" = "user"
  }

  email = each.key
  role  = each.value
}

output "invite_status" {
  value = {
    for email, invite in anthropic_invite.batch :
    email => invite.status
  }
}
