/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package custom_types

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*EmailAddressType)(nil)
)

type EmailAddressType struct {
	basetypes.StringType
}

func (t EmailAddressType) Equal(o attr.Type) bool {
	other, ok := o.(EmailAddressType)
	return ok && t.StringType.Equal(other.StringType)
}

func (t EmailAddressType) String() string {
	return "EmailAddressType"
}

func (t EmailAddressType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := EmailAddressValue{
		StringValue: in,
	}
	return value, nil
}

func (t EmailAddressType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t EmailAddressType) ValueType(_ context.Context) attr.Value {
	return EmailAddressValue{}
}
