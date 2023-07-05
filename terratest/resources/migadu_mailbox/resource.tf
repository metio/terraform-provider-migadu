# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

resource "migadu_mailbox" "mailbox" {
  name        = "Some Name"
  domain_name = var.domain_name
  local_part  = var.local_part
  password    = "secret-token:foobar" // RFC 8959
}
