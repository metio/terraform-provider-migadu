<!--
SPDX-FileCopyrightText: The terraform-provider-migadu Authors
SPDX-License-Identifier: 0BSD
 -->

This file contains migration guidelines for updating the terraform-provider-migadu.

# Migrate to 2023.8.23 or later

The big change here was the implementation of semantic equivalence introduced in [terraform-plugin-framework 1.3](https://github.com/hashicorp/terraform-plugin-framework/issues/70). This made it possible to remove the `_punycode` attributes since we no longer have to differentiate between unicode and ASCII encoded domain names because they are semantically equal. Since removing an attribute is a breaking change anyway, this releases contains another breaking change - the rename of `migadu_rewrite` to `migadu_rewrite_rule` to better reflect what Migadu itself calls these resources. The detailed changes and the proposed action plan is as follows:

## Data Source `migadu_alias`

- The `destinations_punycode` attribute was removed. Use the `destinations` attribute instead. This attribute will contain the destinations as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the destinations in their unicode form.
- The `destinations` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.

## Data Source `migadu_aliases`

- The `aliases[*].destinations_punycode` attribute was removed. Use the `destinations` attribute instead. This attribute will contain the destinations as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the destinations in their unicode form.
- The `aliases[*].destinations` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.

## Resource `migadu_alias`

- The `destinations_punycode` attribute was removed. Put all destinations inside `destinations` attribute instead. You can mix punycode and unicode forms at will and the attribute will retain your formatting.
- The `destinations` attribute is now a set instead of a list. Use the [toset](https://developer.hashicorp.com/terraform/language/functions/toset) function to pass in a list like before.

## Data Source `migadu_mailbox`

- The `delegations_punycode` attribute was removed. Use the `delegations` attribute instead. This attribute will contain the delegations as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the delegations in their unicode form.
- The `delegations` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.
- The `identities_punycode` attribute was removed. Use the `identities` attribute instead. This attribute will contain the identities as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the identities in their unicode form.
- The `identities` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.
- The `recipient_denylist_punycode` attribute was removed. Use the `recipient_denylist` attribute instead. This attribute will contain the recipient denylist as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the recipient denylist in their unicode form.
- The `recipient_denylist` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.
- The `sender_allowlist_punycode` attribute was removed. Use the `sender_allowlist` attribute instead. This attribute will contain the sender allowlist as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the sender allowlist in their unicode form.
- The `sender_allowlist` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.
- The `sender_denylist_punycode` attribute was removed. Use the `sender_denylist` attribute instead. This attribute will contain the sender denylist as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the sender denylist in their unicode form.
- The `sender_denylist` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.

## Data Source `migadu_mailboxes`

- The `mailboxes[*].delegations_punycode` attribute was removed. Use the `delegations` attribute instead. This attribute will contain the delegations as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the delegations in their unicode form.
- The `mailboxes[*].delegations` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.
- The `mailboxes[*].identities_punycode` attribute was removed. Use the `identities` attribute instead. This attribute will contain the identities as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the identities in their unicode form.
- The `mailboxes[*].identities` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.
- The `mailboxes[*].recipient_denylist_punycode` attribute was removed. Use the `recipient_denylist` attribute instead. This attribute will contain the recipient denylist as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the recipient denylist in their unicode form.
- The `mailboxes[*].recipient_denylist` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.
- The `mailboxes[*].sender_allowlist_punycode` attribute was removed. Use the `sender_allowlist` attribute instead. This attribute will contain the sender allowlist as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the sender allowlist in their unicode form.
- The `mailboxes[*].sender_allowlist` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.
- The `mailboxes[*].sender_denylist_punycode` attribute was removed. Use the `sender_denylist` attribute instead. This attribute will contain the sender denylist as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the sender denylist in their unicode form.
- The `mailboxes[*].sender_denylist` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.

## Resource `migadu_mailbox`

- The `delegations_punycode` attribute was removed. Use the `delegations` attribute instead. You can mix punycode and unicode forms at will and the attribute will retain your formatting.
- The `delegations` attribute is now a set instead of a list. Use the [toset](https://developer.hashicorp.com/terraform/language/functions/toset) function to pass in a list like before.
- The `identities_punycode` attribute was removed. Use the `identities` attribute instead. You can mix punycode and unicode forms at will and the attribute will retain your formatting.
- The `identities` attribute is now a set instead of a list. Use the [toset](https://developer.hashicorp.com/terraform/language/functions/toset) function to pass in a list like before.
- The `recipient_denylist_punycode` attribute was removed. Use the `recipient_denylist` attribute instead. You can mix punycode and unicode forms at will and the attribute will retain your formatting.
- The `recipient_denylist` attribute is now a set instead of a list. Use the [toset](https://developer.hashicorp.com/terraform/language/functions/toset) function to pass in a list like before.
- The `sender_allowlist_punycode` attribute was removed. Use the `sender_allowlist` attribute instead. You can mix punycode and unicode forms at will and the attribute will retain your formatting.
- The `sender_allowlist` attribute is now a set instead of a list. Use the [toset](https://developer.hashicorp.com/terraform/language/functions/toset) function to pass in a list like before.
- The `sender_denylist_punycode` attribute was removed. Use the `sender_denylist` attribute instead. You can mix punycode and unicode forms at will and the attribute will retain your formatting.
- The `sender_denylist` attribute is now a set instead of a list. Use the [toset](https://developer.hashicorp.com/terraform/language/functions/toset) function to pass in a list like before.

## Data Source `migadu_rewrite_rule`

- The data source was renamed from `migadu_rewrite` to `migadu_rewrite_rule`
- The `destinations_punycode` attribute was removed. Use the `destinations` attribute instead. This attribute will contain the destinations as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the destinations in their unicode form.
- The `destinations` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.

## Data Source `migadu_rewrite_rules`

- The data source was renamed from `migadu_rewrites` to `migadu_rewrite_rules`
- The `rewrites[*].destinations_punycode` attribute was removed. Use the `destinations` attribute instead. This attribute will contain the destinations as punycode since the Migadu API returns them as such. Please open a ticket in case you need a dedicated attribute containing the destinations in their unicode form.
- The `rewrites[*].destinations` attribute is now a set instead of a list. Use the [tolist](https://developer.hashicorp.com/terraform/language/functions/tolist) function to get a list in case you need one.

## Resource `migadu_rewrite_rule`

- The resource was renamed from `migadu_rewrite` to `migadu_rewrite_rule`
- The `destinations_punycode` attribute was removed. Put all destinations inside `destinations` attribute instead. You can mix punycode and unicode forms at will and the attribute will retain your formatting.
- The `destinations` attribute is now a set instead of a list. Use the [toset](https://developer.hashicorp.com/terraform/language/functions/toset) function to pass in a list like before.
