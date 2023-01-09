# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "id" {
  value = data.migadu_identity.identity.id
}

output "domain_name" {
  value = data.migadu_identity.identity.domain_name
}

output "local_part" {
  value = data.migadu_identity.identity.local_part
}

output "identity" {
  value = data.migadu_identity.identity.identity
}

output "address" {
  value = data.migadu_identity.identity.address
}

output "footer_active" {
  value = data.migadu_identity.identity.footer_active
}

output "footer_html_body" {
  value = data.migadu_identity.identity.footer_html_body
}

output "footer_plain_body" {
  value = data.migadu_identity.identity.footer_plain_body
}

output "may_access_imap" {
  value = data.migadu_identity.identity.may_access_imap
}

output "may_access_manage_sieve" {
  value = data.migadu_identity.identity.may_access_manage_sieve
}

output "may_access_pop3" {
  value = data.migadu_identity.identity.may_access_pop3
}

output "may_receive" {
  value = data.migadu_identity.identity.may_receive
}

output "may_send" {
  value = data.migadu_identity.identity.may_send
}

output "name" {
  value = data.migadu_identity.identity.name
}
