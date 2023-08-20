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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"golang.org/x/net/idna"
)

var (
	_ basetypes.StringTypable = (*DomainNameType)(nil)
	_ xattr.TypeWithValidate  = (*DomainNameType)(nil)
)

type DomainNameType struct {
	basetypes.StringType
}

func (t DomainNameType) Equal(o attr.Type) bool {
	other, ok := o.(DomainNameType)
	return ok && t.StringType.Equal(other.StringType)
}

func (t DomainNameType) String() string {
	return "DomainNameType"
}

func (t DomainNameType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := DomainNameValue{
		StringValue: in,
	}
	return value, nil
}

func (t DomainNameType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t DomainNameType) ValueType(_ context.Context) attr.Value {
	return DomainNameValue{}
}

func (t DomainNameType) Validate(_ context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
	if value.IsNull() || !value.IsKnown() {
		return nil
	}

	var diags diag.Diagnostics
	var valueString string

	if err := value.As(&valueString); err != nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid Terraform Value",
			"An unexpected error occurred while attempting to convert a Terraform value to a string. "+
				"This generally is an issue with the provider schema implementation. "+
				"Please contact the provider developers.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Error: "+err.Error(),
		)
		return diags
	}

	_, err := idna.Lookup.ToASCII(valueString)
	if err != nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid Domain Name String Value",
			"Domain names must be convertible to ASCII.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Given Value: "+valueString+"\n"+
				"Error: "+err.Error(),
		)
		return diags
	}

	return diags
}
