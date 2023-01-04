resource "migadu_identity" "example" {
  domain_name = "example.com"
  local_part  = "some-mailbox"
  identity    = "some-identity"
}

# international domain names are supported
resource "migadu_identity" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-mailbox"
  identity    = "some-identity"
}
