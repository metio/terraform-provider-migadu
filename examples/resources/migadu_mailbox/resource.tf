resource "migadu_mailbox" "example" {
  domain_name = "example.com"
  local_part  = "some-mailbox"
  password    = "Sup3r_s3cr3T"
}

# send invitation to users and let them set password themselves
resource "migadu_mailbox" "invitation" {
  domain_name             = "example.com"
  local_part              = "some-mailbox"
  password_recovery_email = "old@address.example"
}

# international domain names are supported
resource "migadu_mailbox" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-mailbox"
  password    = "Sup3r_s3cr3T"
}
