data "migadu_rewrite_rule" "rewrite" {
  domain_name = "example.com"
  name        = "some-rule"
}

# international domain names are supported
data "migadu_rewrite_rule" "idn" {
  domain_name = "bücher.example"
  name        = "some-rule"
}
