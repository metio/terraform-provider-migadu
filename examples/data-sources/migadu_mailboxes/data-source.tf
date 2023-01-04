data "migadu_mailboxes" "mailboxes" {
  domain_name = "example.com"
}

# international domain names are supported
data "migadu_mailboxes" "idn" {
  domain_name = "bücher.example"
}
