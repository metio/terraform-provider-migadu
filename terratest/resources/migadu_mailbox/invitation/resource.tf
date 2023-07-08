# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

resource "migadu_mailbox" "mailbox" {
  name                    = "Some Name"
  domain_name             = var.domain_name
  local_part              = var.local_part
  password_method         = "invitation"
  password_recovery_email = "someone@example.com"
}
