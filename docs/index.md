---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "migadu Provider"
subcategory: ""
description: |-
  Provider for the Migadu https://www.migadu.com/api/ API. Requires Terraform 1.0 or later.
---

# migadu Provider

Provider for the [Migadu](https://www.migadu.com/api/) API. Requires Terraform 1.0 or later.

## Example Usage

```terraform
provider "migadu" {
  # read all configuration parameters from environment variables
}

# terraform configuration overrides environment variables
provider "migadu" {
  username = "some-name@example.com"
  token    = "your-super-secret-token-that-should-not-be-committed-in-plaintext"
  timeout  = 35
  endpoint = "https://api.migadu.com/v1/"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `endpoint` (String) The API endpoint to use. Can be specified with the `MIGADU_ENDPOINT` environment variable. Defaults to `https://api.migadu.com/v1/`. Take a look at https://www.migadu.com/api/#api-requests for more information.
- `timeout` (Number) The timeout to apply for HTTP requests in seconds. Can be specified with the `MIGADU_TIMEOUT` environment variable. Defaults to `10`.
- `token` (String, Sensitive) The API key to use. Can be specified with the `MIGADU_TOKEN` environment variable. Take a look at https://www.migadu.com/api/#api-keys for more information.
- `username` (String, Sensitive) The username to use. Can be specified with the `MIGADU_USERNAME` environment variable. Take a look at https://www.migadu.com/api/#api-requests for more information.
