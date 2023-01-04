data "migadu_mailbox" "mailbox" {
  domain_name = "example.com"
  local_part  = "some-name"
}

# international domain names are supported
data "migadu_mailbox" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-name"
}
