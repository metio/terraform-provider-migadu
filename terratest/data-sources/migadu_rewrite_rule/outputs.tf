# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

output "id" {
  value = data.migadu_rewrite_rule.rewrite.id
}

output "domain_name" {
  value = data.migadu_rewrite_rule.rewrite.domain_name
}

output "name" {
  value = data.migadu_rewrite_rule.rewrite.name
}

output "local_part_rule" {
  value = data.migadu_rewrite_rule.rewrite.local_part_rule
}

output "order_num" {
  value = data.migadu_rewrite_rule.rewrite.order_num
}

output "destinations" {
  value = data.migadu_rewrite_rule.rewrite.destinations
}
