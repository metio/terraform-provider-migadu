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
)

var (
	_ basetypes.SetValuable                   = (*EmailAddressSetValue)(nil)
	_ basetypes.SetValuableWithSemanticEquals = (*EmailAddressSetValue)(nil)
)

type EmailAddressSetValue struct {
	basetypes.SetValue
}

func (v EmailAddressSetValue) Type(_ context.Context) attr.Type {
	return EmailAddressSetType{
		SetType: basetypes.SetType{ElemType: EmailAddressType{}},
	}
}

func (v EmailAddressSetValue) Equal(o attr.Value) bool {
	other, ok := o.(EmailAddressSetValue)

	if !ok {
		return false
	}

	return v.SetValue.Equal(other.SetValue)
}

func (v EmailAddressSetValue) SetSemanticEquals(ctx context.Context, newValuable basetypes.SetValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(EmailAddressSetValue)

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

	if !v.ElementType(ctx).Equal(newValue.ElementType(ctx)) {
		return false, diags
	}

	if len(v.Elements()) != len(newValue.Elements()) {
		return false, diags
	}

	for _, elem := range v.Elements() {
		if !newValue.contains(ctx, elem) {
			return false, diags
		}
	}

	return true, diags
}

func (v EmailAddressSetValue) contains(ctx context.Context, other attr.Value) bool {
	otherEmail, ok := other.(EmailAddressValue)

	if ok {
		for _, elem := range v.Elements() {
			if email, ok := elem.(EmailAddressValue); ok {
				if equal, _ := email.StringSemanticEquals(ctx, otherEmail); equal {
					return true
				}
			}
		}
	}

	return false
}
