/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package custom_types

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

// NewEmailAddressValue creates an email with a known value.
func NewEmailAddressValue(value string) EmailAddressValue {
	return EmailAddressValue{
		StringValue: basetypes.NewStringValue(value),
	}
}
