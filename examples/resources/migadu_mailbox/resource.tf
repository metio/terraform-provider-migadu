resource "migadu_mailbox" "example" {
  domain_name = "example.com"
  local_part  = "some-mailbox"
  password    = "Sup3r_s3cr3T"
}

# international domain names
resource "migadu_mailbox" "idn" {
  domain_name = "bücher.example"
  local_part  = "some-mailbox"
  password    = "Sup3r_s3cr3T"
}
