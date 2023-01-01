package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"golang.org/x/net/idna"
	"strings"
)

func ConvertEmailsToUnicode(emails []string, diag *diag.Diagnostics) []string {
	return convertEmailsWith(emails, diag, idna.ToUnicode)
}

func ConvertEmailsToASCII(emails []string, diag *diag.Diagnostics) []string {
	return convertEmailsWith(emails, diag, idna.ToASCII)
}

func convertEmailsWith(originals []string, diag *diag.Diagnostics, converter func(string) (string, error)) []string {
	var modified []string
	for _, email := range originals {
		parts := strings.Split(email, "@")
		converted, err := converter(parts[1])
		if err == nil {
			modified = append(modified, fmt.Sprintf("%s@%s", parts[0], converted))
		} else {
			diag.AddError(
				"Error converting email",
				fmt.Sprintf("Could not convert %s: %v", email, err),
			)
		}
	}
	return modified
}
