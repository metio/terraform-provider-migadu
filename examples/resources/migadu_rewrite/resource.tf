resource "migadu_rewrite" "example" {
  domain_name     = "example.com"
  name            = "security-mails"
  local_part_rule = "sec-*"

  destinations = [
    "first@example.com",
    "second@example.com",
  ]
}

# international domain names are supported
resource "migadu_rewrite" "idn" {
  domain_name     = "bücher.example"
  name            = "security-mails"
  local_part_rule = "sec-*"

  destinations = [
    "first@ücher.example",
    "second@ücher.example",
  ]
}
