# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

output "id" {
  value = data.migadu_rewrite.rewrite.id
}

output "domain_name" {
  value = data.migadu_rewrite.rewrite.domain_name
}

output "name" {
  value = data.migadu_rewrite.rewrite.name
}

output "local_part_rule" {
  value = data.migadu_rewrite.rewrite.local_part_rule
}

output "order_num" {
  value = data.migadu_rewrite.rewrite.order_num
}

output "destinations" {
  value = data.migadu_rewrite.rewrite.destinations
}

output "destinations_punycode" {
  value = data.migadu_rewrite.rewrite.destinations_punycode
}
