data "migadu_rewrite" "rewrite" {
  domain_name = "example.com"
  name        = "some-rule"
}

# international domain names are supported
data "migadu_rewrite" "idn" {
  domain_name = "bücher.example"
  name        = "some-rule"
}
