package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
)

func CreateIdentityID(localPart types.String, domainName custom_types.DomainNameValue, identity types.String) string {
	return CreateIdentityIDString(localPart.ValueString(), domainName.ValueString(), identity.ValueString())
}

func CreateIdentityIDString(localPart, domainName, identity string) string {
	return fmt.Sprintf("%s@%s/%s", localPart, domainName, identity)
}

func IdentityCreateError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Creating Identity",
		standardAPIErrorDetail(err),
	)
}

func IdentityReadError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Reading Identity",
		standardAPIErrorDetail(err),
	)
}

func IdentityUpdateError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Updating Identity",
		standardAPIErrorDetail(err),
	)
}

func IdentityDeleteError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Deleting Identity",
		standardAPIErrorDetail(err),
	)
}

func IdentityImportError(id string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Importing Identity",
		standardImportErrorDetail("local_part@domain_name/identity", id),
	)
}
