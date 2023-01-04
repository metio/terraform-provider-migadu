data "migadu_aliases" "aliases" {
  domain_name = "example.com"
}

# international domain names are supported
data "migadu_aliases" "idn" {
  domain_name = "bücher.example"
}
