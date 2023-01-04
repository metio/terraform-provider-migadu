data "migadu_alias" "alias" {
  domain_name = "example.com"
  local_part  = "some-name"
}

# international domain names are supported
data "migadu_alias" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-name"
}
