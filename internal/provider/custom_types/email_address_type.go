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
	"strings"
)

var (
	_ basetypes.StringTypable = (*EmailAddressType)(nil)
	_ xattr.TypeWithValidate  = (*EmailAddressType)(nil)
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

func (t EmailAddressType) Validate(_ context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
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

	if len(valueString) == 0 {
		return nil
	}

	valueParts := strings.Split(valueString, "@")

	if len(valueParts) != 2 || len(valueParts[0]) == 0 || len(valueParts[1]) == 0 {
		diags.AddAttributeError(
			valuePath,
			"Invalid Email Address String Value",
			"An email must match the format 'local_part@domain'.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Given Value: "+valueString,
		)
		return diags
	}

	return diags
}
