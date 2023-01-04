data "migadu_identity" "identity" {
  domain_name = "example.com"
  local_part  = "mailbox"
  identity    = "some-identity"
}

# international domain names are supported
data "migadu_identity" "idn" {
  domain_name = "bücher.example"
  local_part  = "mailbox"
  identity    = "some-identity"
}
