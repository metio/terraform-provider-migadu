# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "id" {
  value = data.migadu_mailbox.mailbox.id
}

output "domain_name" {
  value = data.migadu_mailbox.mailbox.domain_name
}

output "local_part" {
  value = data.migadu_mailbox.mailbox.local_part
}

output "address" {
  value = data.migadu_mailbox.mailbox.address
}

output "name" {
  value = data.migadu_mailbox.mailbox.name
}

output "delegations" {
  value = data.migadu_mailbox.mailbox.delegations
}

output "delegations_punycode" {
  value = data.migadu_mailbox.mailbox.delegations_punycode
}

output "auto_respond_active" {
  value = data.migadu_mailbox.mailbox.auto_respond_active
}

output "expirable" {
  value = data.migadu_mailbox.mailbox.expirable
}

output "expires_on" {
  value = data.migadu_mailbox.mailbox.expires_on
}

output "is_internal" {
  value = data.migadu_mailbox.mailbox.is_internal
}

output "remove_upon_expiry" {
  value = data.migadu_mailbox.mailbox.remove_upon_expiry
}
