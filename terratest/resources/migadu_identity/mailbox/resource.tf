# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

resource "migadu_identity" "identity" {
  name         = "Some Name"
  domain_name  = var.domain_name
  local_part   = var.local_part
  identity     = var.identity
  password_use = "mailbox"
}
