/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package idn

import (
	"fmt"
	"golang.org/x/net/idna"
	"strings"
)

// ConvertEmailsToUnicode converts the domain name of email addresses to their unicode representation
func ConvertEmailsToUnicode(emails []string) ([]string, error) {
	return convertEmailsWith(emails, idna.ToUnicode)
}

// ConvertEmailsToASCII converts the domain name of email addresses to their punycode representation
func ConvertEmailsToASCII(emails []string) ([]string, error) {
	return convertEmailsWith(emails, idna.ToASCII)
}

func convertEmailsWith(originals []string, converter func(string) (string, error)) ([]string, error) {
	var modified []string
	for _, email := range originals {
		parts := strings.Split(email, "@")
		converted, err := converter(parts[1])
		if err != nil {
			return nil, err
		}
		modified = append(modified, fmt.Sprintf("%s@%s", parts[0], converted))
	}
	return modified, nil
}
