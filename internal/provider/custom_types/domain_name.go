package custom_types

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

// NewDomainNameValue creates an email with a known value.
func NewDomainNameValue(value string) DomainNameValue {
	return DomainNameValue{
		StringValue: basetypes.NewStringValue(value),
	}
}
