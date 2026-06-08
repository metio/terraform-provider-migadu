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

# override the default rate limit (defaults to 60 requests per 2m)
provider "migadu" {
  username      = "some-name@example.com"
  token         = "your-super-secret-token-that-should-not-be-committed-in-plaintext"
  rate_limit    = 30
  rate_interval = "1m"
}

# disable client-side rate limiting
provider "migadu" {
  username   = "some-name@example.com"
  token      = "your-super-secret-token-that-should-not-be-committed-in-plaintext"
  rate_limit = 0
}
