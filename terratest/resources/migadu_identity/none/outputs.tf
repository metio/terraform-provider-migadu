# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

output "id" {
  value = migadu_identity.identity.id
}

output "domain_name" {
  value = migadu_identity.identity.domain_name
}

output "local_part" {
  value = migadu_identity.identity.local_part
}

output "address" {
  value = migadu_identity.identity.address
}

output "name" {
  value = migadu_identity.identity.name
}
