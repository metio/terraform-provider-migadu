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
	"github.com/metio/terraform-provider-migadu/migadu/idn"
	"strings"
)

var (
	_ basetypes.StringValuable                   = (*EmailAddressValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*EmailAddressValue)(nil)
)

type EmailAddressValue struct {
	basetypes.StringValue
}

func (v EmailAddressValue) Type(_ context.Context) attr.Type {
	return EmailAddressType{}
}

func (v EmailAddressValue) Equal(o attr.Value) bool {
	other, ok := o.(EmailAddressValue)
	return ok && v.StringValue.Equal(other.StringValue)
}

func (v EmailAddressValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(EmailAddressValue)

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

	priorEmail, err := normalizeEmail(v.StringValue.ValueString())
	if err != nil {
		diags.AddError(
			"Semantic Equality Check Error",
			"An error occurred while normalizing an email address. "+
				"Please report this to the provider developers.\n\n"+
				"Given Value: "+v.StringValue.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
		return false, diags
	}

	newEmail, err := normalizeEmail(newValue.ValueString())
	if err != nil {
		diags.AddError(
			"Semantic Equality Check Error",
			"An error occurred while normalizing an email address. "+
				"Please report this to the provider developers.\n\n"+
				"Given Value: "+newValue.ValueString()+"\n"+
				"Error: "+err.Error(),
		)
		return false, diags
	}

	return priorEmail == newEmail, diags
}

func normalizeEmail(email string) (string, error) {
	normalized := email
	normalized = strings.TrimSpace(normalized)
	normalized = strings.ToLower(normalized)
	return idn.ConvertEmailToASCII(normalized)
}
