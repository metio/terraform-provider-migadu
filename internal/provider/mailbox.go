/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
)

func CreateMailboxID(localPart types.String, domainName custom_types.DomainNameValue) string {
	return CreateMailboxIDString(localPart.ValueString(), domainName.ValueString())
}

func CreateMailboxIDString(localPart, domainName string) string {
	return fmt.Sprintf("%s@%s", localPart, domainName)
}

func MailboxCreateError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Creating Mailbox",
		standardAPIErrorDetail(err),
	)
}

func MailboxReadError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Reading Mailbox",
		standardAPIErrorDetail(err),
	)
}

func MailboxUpdateError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Updating Mailbox",
		standardAPIErrorDetail(err),
	)
}

func MailboxDeleteError(err error) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Deleting Mailbox",
		standardAPIErrorDetail(err),
	)
}

func MailboxImportError(id string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Error Importing Mailbox",
		standardImportErrorDetail("local_part@domain_name", id),
	)
}
