# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

output "id" {
  value = migadu_mailbox.mailbox.id
}

output "domain_name" {
  value = migadu_mailbox.mailbox.domain_name
}

output "local_part" {
  value = migadu_mailbox.mailbox.local_part
}

output "address" {
  value = migadu_mailbox.mailbox.address
}

output "name" {
  value = migadu_mailbox.mailbox.name
}

output "delegations" {
  value = migadu_mailbox.mailbox.delegations
}

output "identities" {
  value = migadu_mailbox.mailbox.identities
}

output "auto_respond_active" {
  value = migadu_mailbox.mailbox.auto_respond_active
}

output "expirable" {
  value = migadu_mailbox.mailbox.expirable
}

output "expires_on" {
  value = migadu_mailbox.mailbox.expires_on
}

output "is_internal" {
  value = migadu_mailbox.mailbox.is_internal
}

output "remove_upon_expiry" {
  value = migadu_mailbox.mailbox.remove_upon_expiry
}

output "password_recovery_email" {
  value = migadu_mailbox.mailbox.password_recovery_email
}
