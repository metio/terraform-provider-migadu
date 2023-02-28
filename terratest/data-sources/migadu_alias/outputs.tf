# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

output "id" {
  value = data.migadu_alias.alias.id
}

output "domain_name" {
  value = data.migadu_alias.alias.domain_name
}

output "local_part" {
  value = data.migadu_alias.alias.local_part
}

output "address" {
  value = data.migadu_alias.alias.address
}

output "destinations" {
  value = data.migadu_alias.alias.destinations
}

output "destinations_punycode" {
  value = data.migadu_alias.alias.destinations_punycode
}

output "expirable" {
  value = data.migadu_alias.alias.expirable
}

output "expires_on" {
  value = data.migadu_alias.alias.expires_on
}

output "is_internal" {
  value = data.migadu_alias.alias.is_internal
}

output "remove_upon_expiry" {
  value = data.migadu_alias.alias.remove_upon_expiry
}
