# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

data "migadu_rewrite_rule" "rewrite" {
  domain_name = var.domain_name
  name        = var.name
}
