resource "migadu_identity" "example" {
  domain_name = "example.com"
  local_part  = "some-mailbox"
  identity    = "some-identity"
  name        = "Some Name"
}

# international domain names are supported
resource "migadu_identity" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-mailbox"
  identity    = "some-identity"
  name        = "Some Name"
}

# Sender persona, not used for authentication
resource "migadu_identity" "sender_persona" {
  domain_name  = "example.com"
  local_part   = "some-mailbox"
  identity     = "some-identity"
  name         = "Some Name"
  password_use = "none" # use mailbox user/password
}

# alternative address linked to the same mailbox
resource "migadu_identity" "mailbox" {
  domain_name  = "example.com"
  local_part   = "some-mailbox"
  identity     = "some-identity"
  name         = "Some Name"
  password_use = "mailbox" # use identity user and mailbox password
}

# application specific password
resource "migadu_identity" "custom" {
  domain_name  = "example.com"
  local_part   = "some-mailbox"
  identity     = "some-identity"
  name         = "Some Name"
  password_use = "custom" # use identity user/password
  password     = "Sup3r_s3cr3T"
}
