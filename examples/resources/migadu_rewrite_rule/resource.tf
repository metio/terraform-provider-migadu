resource "migadu_rewrite_rule" "example" {
  domain_name     = "example.com"
  name            = "security-mails"
  local_part_rule = "sec-*"

  destinations = [
    "first@example.com",
    "second@example.com",
  ]
}

# international domain names are supported
resource "migadu_rewrite_rule" "idn" {
  domain_name     = "bücher.example"
  name            = "security-mails"
  local_part_rule = "sec-*"

  destinations = [
    "first@bücher.example",
    "second@bücher.example",
  ]
}
