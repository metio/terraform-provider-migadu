package provider

import (
	"fmt"
)

func standardAPIErrorDetail(err error) string {
	return "While calling the API, an unexpected error was returned in the response. " +
		"Please contact the provider developer if you are unsure how to resolve the error.\n\n" +
		"Error: " + err.Error()
}

func standardImportErrorDetail(format string, id string) string {
	return fmt.Sprintf("Expected import identifier with format: '%s' Got: '%s'", format, id)
}
