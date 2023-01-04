data "migadu_identities" "identities" {
  domain_name = "example.com"
  local_part  = "some-name"
}

# international domain names are supported
data "migadu_identities" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-name"
}
