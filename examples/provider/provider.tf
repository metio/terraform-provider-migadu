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
