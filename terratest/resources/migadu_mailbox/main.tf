# SPDX-FileCopyrightText: The terraform-provider-migadu Authors
# SPDX-License-Identifier: 0BSD

terraform {
  required_providers {
    migadu = {
      source  = "localhost/metio/migadu"
      version = "9999.99.99"
    }
  }
}

provider "migadu" {
  username = "terratest"
  token    = "secret-token:foobar" // RFC 8959
  endpoint = var.endpoint
}
