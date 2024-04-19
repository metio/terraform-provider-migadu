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
	_ basetypes.SetTypable = (*EmailAddressSetType)(nil)
)

type EmailAddressSetType struct {
	basetypes.SetType
}

func (t EmailAddressSetType) Equal(o attr.Type) bool {
	other, ok := o.(EmailAddressSetType)

	if !ok {
		return false
	}

	return t.SetType.Equal(other.SetType)
}

func (t EmailAddressSetType) String() string {
	return fmt.Sprintf("EmailAddressSetType(%s)", t.SetType.String())
}

func (t EmailAddressSetType) ValueFromSet(_ context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	value := EmailAddressSetValue{
		SetValue: in,
	}

	return value, diags
}

func (t EmailAddressSetType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	setValue, ok := attrValue.(basetypes.SetValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	setValuable, diags := t.ValueFromSet(ctx, setValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting SetValue to SetValuable: %v", diags)
	}

	return setValuable, nil
}

func (t EmailAddressSetType) ValueType(_ context.Context) attr.Value {
	return EmailAddressSetValue{}
}
