package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
)

func CreateAliasID(localPart types.String, domainName custom_types.DomainNameValue) string {
	return CreateAliasIDString(localPart.ValueString(), domainName.ValueString())
}

func CreateAliasIDString(localPart string, domainName string) string {
	return fmt.Sprintf("%s@%s", localPart, domainName)
}

func AliasCreateError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Creating Alias",
		standardAPIErrorDetail(err),
	)
}

func AliasReadError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Reading Alias",
		standardAPIErrorDetail(err),
	)
}

func AliasUpdateError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Updating Alias",
		standardAPIErrorDetail(err),
	)
}

func AliasDeleteError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Deleting Alias",
		standardAPIErrorDetail(err),
	)
}

func AliasImportError(id string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Importing Alias",
		standardImportErrorDetail("local_part@domain_name", id),
	)
}
