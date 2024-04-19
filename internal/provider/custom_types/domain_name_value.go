/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package custom_types

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"golang.org/x/net/idna"
	"strings"
)

var (
	_ basetypes.StringValuable                   = (*DomainNameValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*DomainNameValue)(nil)
	_ xattr.ValidateableAttribute                = (*DomainNameValue)(nil)
)

// NewDomainNameValue creates a domain name with a known value.
func NewDomainNameValue(value string) DomainNameValue {
	return DomainNameValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

type DomainNameValue struct {
	basetypes.StringValue
}

func (v DomainNameValue) Type(_ context.Context) attr.Type {
	return DomainNameType{}
}

func (v DomainNameValue) Equal(o attr.Value) bool {
	other, ok := o.(DomainNameValue)
	return ok && v.StringValue.Equal(other.StringValue)
}

func (v DomainNameValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(DomainNameValue)

	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)
		return false, diags
	}

	priorDomain, err := normalizeDomain(v.StringValue.ValueString())
	if err != nil {
		diags.AddError(
			"Semantic Equality Check Error",
			"An error occurred while normalizing a domain name. "+
				"Please report this to the provider developers.\n\n"+
				"Given Value: "+v.StringValue.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
		return false, diags
	}

	newDomain, err := normalizeDomain(newValue.ValueString())
	if err != nil {
		diags.AddError(
			"Semantic Equality Check Error",
			"An error occurred while normalizing a domain name. "+
				"Please report this to the provider developers.\n\n"+
				"Given Value: "+newValue.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
		return false, diags
	}

	return priorDomain == newDomain, diags
}

func normalizeDomain(domain string) (string, error) {
	normalized := domain
	normalized = strings.TrimSpace(normalized)
	normalized = strings.ToLower(normalized)
	return idna.ToASCII(normalized)
}

func (v DomainNameValue) ValidateAttribute(_ context.Context, request xattr.ValidateAttributeRequest, response *xattr.ValidateAttributeResponse) {
	if v.IsNull() || v.IsUnknown() {
		return
	}

	_, err := idna.Lookup.ToASCII(v.ValueString())
	if err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Domain Name String Value",
			"Domain names must be convertible to ASCII.\n\n"+
				"Path: "+request.Path.String()+"\n"+
				"Given Value: "+v.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
		return
	}
}
