/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/metio/terraform-provider-migadu/migadu/idn"
)

func ConvertEmailsToUnicode(emails []string, diag *diag.Diagnostics) []string {
	converted, err := idn.ConvertEmailsToUnicode(emails)
	if err != nil {
		diag.AddError("Error converting email", err.Error())
		return []string{}
	}
	return converted
}

func ConvertEmailsToASCII(emails []string, diag *diag.Diagnostics) []string {
	converted, err := idn.ConvertEmailsToASCII(emails)
	if err != nil {
		diag.AddError("Error converting email", err.Error())
		return []string{}
	}
	return converted
}
