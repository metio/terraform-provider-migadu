# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

data "migadu_mailbox" "mailbox" {
  domain_name = var.domain_name
  local_part  = var.local_part
}
