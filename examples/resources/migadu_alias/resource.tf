resource "migadu_alias" "example" {
  domain_name = "example.com"
  local_part  = "some-name"

  destinations = [
    "first@example.com",
    "second@example.com",
  ]
}

# international domain names are supported
resource "migadu_alias" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-name"

  destinations = [
    "first@bücher.example",
    "second@bücher.example",
  ]
}
