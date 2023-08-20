package custom_types

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewEmailAddressSetValueFrom(ctx context.Context, emails []string) (EmailAddressSetValue, diag.Diagnostics) {
	setValue, diagnostics := basetypes.NewSetValueFrom(ctx, EmailAddressType{}, emails)
	return EmailAddressSetValue{
		SetValue: setValue,
	}, diagnostics
}

func NewEmailAddressSetNull() EmailAddressSetValue {
	setValue := basetypes.NewSetNull(EmailAddressType{})
	return EmailAddressSetValue{
		SetValue: setValue,
	}
}
