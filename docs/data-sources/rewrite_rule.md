---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "migadu_rewrite_rule Data Source - terraform-provider-migadu"
subcategory: ""
description: |-
  Get information about a single rewrite rule.
---

# migadu_rewrite_rule (Data Source)

Get information about a single rewrite rule.

## Example Usage

```terraform
data "migadu_rewrite_rule" "rewrite" {
  domain_name = "example.com"
  name        = "some-rule"
}

# international domain names are supported
data "migadu_rewrite_rule" "idn" {
  domain_name = "bücher.example"
  name        = "some-rule"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_name` (String) The domain of the rewrite rule.
- `name` (String) The name (slug) of the rewrite rule.

### Read-Only

- `destinations` (Set of String) The destinations of the rewrite rule.
- `id` (String) Contains the value `domain_name/name`.
- `local_part_rule` (String) The local part expression of the rewrite rule
- `order_num` (Number) The order number of the rewrite rule.