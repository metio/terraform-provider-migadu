/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package custom_validators

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = passwordUseValidator{}

type passwordUseValidator struct{}

func (v passwordUseValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v passwordUseValidator) MarkdownDescription(_ context.Context) string {
	return "password_use attribute is invalid"
}

func (v passwordUseValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	stringvalidator.OneOf("none", "mailbox", "custom").ValidateString(ctx, request, response)

	if value.ValueString() == "none" || value.ValueString() == "mailbox" {
		stringvalidator.ConflictsWith(path.MatchRoot("password")).ValidateString(ctx, request, response)
	} else if value.ValueString() == "custom" {
		stringvalidator.AlsoRequires(path.MatchRoot("password")).ValidateString(ctx, request, response)
	}
}

// PasswordUse validates the `password_use` attribute of an identity resource
func PasswordUse() validator.String {
	return passwordUseValidator{}
}
