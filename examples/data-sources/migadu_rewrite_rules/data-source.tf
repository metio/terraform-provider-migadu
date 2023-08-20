data "migadu_rewrite_rules" "rewrites" {
  domain_name = "example.com"
}

# international domain names are supported
data "migadu_rewrite_rules" "idn" {
  domain_name = "bücher.example"
}
