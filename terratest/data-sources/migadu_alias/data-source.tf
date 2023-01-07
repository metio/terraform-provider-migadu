# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

data "migadu_alias" "alias" {
  domain_name = var.domain_name
  local_part  = var.local_part
}
