package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
)

func CreateRewriteRuleID(domainName custom_types.DomainNameValue, name types.String) string {
	return CreateRewriteRuleIDString(domainName.ValueString(), name.ValueString())
}

func CreateRewriteRuleIDString(domainName, name string) string {
	return fmt.Sprintf("%s/%s", domainName, name)
}

func RewriteRuleCreateError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Creating RewriteRule Rule",
		standardAPIErrorDetail(err),
	)
}

func RewriteRuleReadError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Reading RewriteRule Rule",
		standardAPIErrorDetail(err),
	)
}

func RewriteRuleUpdateError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Updating RewriteRule Rule",
		standardAPIErrorDetail(err),
	)
}

func RewriteRuleDeleteError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Deleting RewriteRule Rule",
		standardAPIErrorDetail(err),
	)
}

func RewriteRuleImportError(id string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Importing RewriteRule Rule",
		standardImportErrorDetail("domain_name/name", id),
	)
}
