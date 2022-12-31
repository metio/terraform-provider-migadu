resource "migadu_alias" "example" {
  domain_name = "example.com"
  local_part  = "some-name"

  destinations = [
    "first@example.com",
    "second@example.com",
  ]
}

# international domain names
resource "migadu_alias" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-name"

  # notice that we are using 'destinations_idn' here
  destinations_idn = [
    "first@bücher.example",
    "second@bücher.example",
  ]
}
