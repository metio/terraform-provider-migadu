package custom_types

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.SetTypable   = (*EmailAddressSetType)(nil)
	_ xattr.TypeWithValidate = (*EmailAddressSetType)(nil)
)

var emailAddressSetType = types.SetType{
	ElemType: EmailAddressType{},
}

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

func (t EmailAddressSetType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	if !in.Type().Is(tftypes.Set{}) {
		err := fmt.Errorf("expected Set value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Set Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var elems []tftypes.Value

	if err := in.As(&elems); err != nil {
		diags.AddAttributeError(
			path,
			"Set Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

	validatableType, isValidatable := t.ElementType().(xattr.TypeWithValidate)

	for indexOuter, elemOuter := range elems {
		if !elemOuter.IsFullyKnown() {
			continue
		}

		outerValue, err := t.ElementType().ValueFromTerraform(ctx, elemOuter)
		if err != nil {
			diags.AddAttributeError(
				path,
				"Set Type Validation Error",
				"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)
			return diags
		}
		outerEmail := outerValue.(EmailAddressValue)

		// Validate the element first
		if isValidatable {
			diags = append(diags, validatableType.Validate(ctx, elemOuter, path.AtSetValue(outerValue))...)
		}

		// Then check for duplicates
		for indexInner := indexOuter + 1; indexInner < len(elems); indexInner++ {
			innerValue, err := t.ElementType().ValueFromTerraform(ctx, elems[indexInner])
			if err != nil {
				diags.AddAttributeError(
					path,
					"Set Type Validation Error",
					"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
				)
				return diags
			}
			innerEmail := innerValue.(EmailAddressValue)

			if equal, _ := innerEmail.StringSemanticEquals(ctx, outerEmail); !equal {
				continue
			}

			diags.AddAttributeError(
				path,
				"Duplicate Set Element",
				fmt.Sprintf("This attribute contains duplicate values of: %s", innerValue),
			)
		}
	}

	return diags
}
