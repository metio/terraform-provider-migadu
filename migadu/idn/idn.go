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

// ConvertEmailToASCII converts the domain name of an email address to its punycode representation
func ConvertEmailToASCII(email string) (string, error) {
	return convertEmailWith(email, idna.ToASCII)
}

// ConvertEmailsToASCII converts the domain name of email addresses to their punycode representation
func ConvertEmailsToASCII(emails []string) ([]string, error) {
	return convertEmailsWith(emails, idna.ToASCII)
}

func convertEmailsWith(originals []string, converter func(string) (string, error)) ([]string, error) {
	var modified []string
	for _, email := range originals {
		converted, err := convertEmailWith(email, converter)
		if err != nil {
			return nil, err
		}
		modified = append(modified, converted)
	}
	return modified, nil
}

func convertEmailWith(email string, converter func(string) (string, error)) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email, nil
	}
	converted, err := converter(parts[1])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s@%s", parts[0], converted), nil

}
