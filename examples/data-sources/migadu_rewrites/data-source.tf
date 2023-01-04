data "migadu_rewrites" "rewrites" {
  domain_name = "example.com"
}

# international domain names are supported
data "migadu_rewrites" "idn" {
  domain_name = "bücher.example"
}
