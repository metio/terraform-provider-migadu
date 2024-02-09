/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/migadu-client.go/model"
	"github.com/metio/migadu-client.go/simulator"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
)

func TestMailboxResource_API_Success_Using_Password(t *testing.T) {
	testCases := map[string]ResourceTestCase[model.Mailbox]{
		"change-name": {
			Create: ResourceTestStep[model.Mailbox]{
				Send: model.Mailbox{
					LocalPart:  "test",
					DomainName: "example.com",
					Name:       "Some Name",
					Password:   "secret",
				},
				Want: model.Mailbox{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
					Password:   "secret",
				},
			},
			Update: ResourceTestStep[model.Mailbox]{
				Send: model.Mailbox{
					LocalPart:  "test",
					DomainName: "example.com",
					Name:       "Different Name",
					Password:   "secret",
				},
				Want: model.Mailbox{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Different Name",
					Password:   "secret",
				},
			},
		},
		"change-idna-domain": {
			Create: ResourceTestStep[model.Mailbox]{
				Send: model.Mailbox{
					LocalPart:  "test",
					DomainName: "hoß.de",
					Name:       "Some Name",
					Password:   "secret",
				},
				Want: model.Mailbox{
					LocalPart:  "test",
					DomainName: "hoß.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
					Password:   "secret",
				},
			},
			Update: ResourceTestStep[model.Mailbox]{
				Send: model.Mailbox{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Name:       "Some Name",
					Password:   "secret",
				},
				Want: model.Mailbox{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
					Password:   "secret",
				},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								local_part  = "%s"
								domain_name = "%s"
								password    = "%s"
								name        = "%s"
							}
						`, testCase.Create.Send.LocalPart, testCase.Create.Send.DomainName, testCase.Create.Send.Password, testCase.Create.Send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", testCase.Create.Want.LocalPart, testCase.Create.Want.DomainName)),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", testCase.Create.Want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", testCase.Create.Want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", testCase.Create.Want.Name),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", testCase.Create.Want.Password),
						),
					},
					{
						ResourceName:      "migadu_mailbox.test",
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateVerifyIgnore: []string{
							"password", // Migadu API does not allow reading passwords
						},
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								local_part  = "%s"
								domain_name = "%s"
								password    = "%s"
								name        = "%s"
							}
						`, testCase.Update.Send.LocalPart, testCase.Update.Send.DomainName, testCase.Update.Send.Password, testCase.Update.Send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", testCase.Update.Want.LocalPart, testCase.Update.Want.DomainName)),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", testCase.Update.Want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", testCase.Update.Want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", testCase.Update.Want.Name),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", testCase.Update.Want.Password),
						),
					},
				},
			})
		})
	}
}

func TestMailboxResource_API_Success_Using_RecoveryEmail(t *testing.T) {
	testCases := map[string]ResourceTestCase[model.Mailbox]{
		"change-name": {
			Create: ResourceTestStep[model.Mailbox]{
				Send: model.Mailbox{
					LocalPart:             "test",
					DomainName:            "example.com",
					Name:                  "Some Name",
					PasswordRecoveryEmail: "someone@example.com",
				},
				Want: model.Mailbox{
					LocalPart:             "test",
					DomainName:            "example.com",
					Address:               "test@example.com",
					Name:                  "Some Name",
					PasswordMethod:        "invitation",
					PasswordRecoveryEmail: "someone@example.com",
				},
			},
			Update: ResourceTestStep[model.Mailbox]{
				Send: model.Mailbox{
					LocalPart:             "test",
					DomainName:            "example.com",
					Name:                  "Different Name",
					PasswordRecoveryEmail: "someone@example.com",
				},
				Want: model.Mailbox{
					LocalPart:             "test",
					DomainName:            "example.com",
					Address:               "test@example.com",
					Name:                  "Different Name",
					PasswordMethod:        "invitation",
					PasswordRecoveryEmail: "someone@example.com",
				},
			},
		},
		"change-idna-domain": {
			Create: ResourceTestStep[model.Mailbox]{
				Send: model.Mailbox{
					LocalPart:             "test",
					DomainName:            "hoß.de",
					Name:                  "Some Name",
					PasswordRecoveryEmail: "someone@hoß.de",
				},
				Want: model.Mailbox{
					LocalPart:             "test",
					DomainName:            "hoß.de",
					Address:               "test@xn--ho-hia.de",
					Name:                  "Some Name",
					PasswordMethod:        "invitation",
					PasswordRecoveryEmail: "someone@hoß.de",
				},
			},
			Update: ResourceTestStep[model.Mailbox]{
				Send: model.Mailbox{
					LocalPart:             "test",
					DomainName:            "xn--ho-hia.de",
					Name:                  "Some Name",
					PasswordRecoveryEmail: "someone@hoß.de",
				},
				Want: model.Mailbox{
					LocalPart:             "test",
					DomainName:            "xn--ho-hia.de",
					Address:               "test@xn--ho-hia.de",
					Name:                  "Some Name",
					PasswordMethod:        "invitation",
					PasswordRecoveryEmail: "someone@hoß.de",
				},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								local_part              = "%s"
								domain_name             = "%s"
								password_method         = "invitation"
								password_recovery_email = "%s"
								name                    = "%s"
							}
						`, testCase.Create.Send.LocalPart, testCase.Create.Send.DomainName, testCase.Create.Send.PasswordRecoveryEmail, testCase.Create.Send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", testCase.Create.Want.LocalPart, testCase.Create.Want.DomainName)),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", testCase.Create.Want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", testCase.Create.Want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", testCase.Create.Want.Name),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", ""),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password_recovery_email", testCase.Create.Want.PasswordRecoveryEmail),
						),
					},
					{
						ResourceName:      "migadu_mailbox.test",
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateVerifyIgnore: []string{
							"password", // Migadu API does not allow reading password
						},
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								local_part              = "%s"
								domain_name             = "%s"
								password_recovery_email = "%s"
								password_method         = "invitation"
								name                    = "%s"
							}
						`, testCase.Update.Send.LocalPart, testCase.Update.Send.DomainName, testCase.Update.Send.PasswordRecoveryEmail, testCase.Update.Send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", testCase.Update.Want.LocalPart, testCase.Update.Want.DomainName)),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", testCase.Update.Want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", testCase.Update.Want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", testCase.Update.Want.Name),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", ""),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password_recovery_email", testCase.Update.Want.PasswordRecoveryEmail),
						),
					},
				},
			})
		})
	}
}

func TestMailboxResource_API_Errors(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-400": {
			StatusCode: http.StatusBadRequest,
			ErrorRegex: "CreateMailbox: status: 400",
		},
		"error-409": {
			StatusCode: http.StatusConflict,
			ErrorRegex: "CreateMailbox: status: 409",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "CreateMailbox: status: 500",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{StatusCode: testCase.StatusCode}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + `
							resource "migadu_mailbox" "test" {
								name        = "Some Name"
								local_part  = "test"
								domain_name = "example.com"
								password    = "secret"
							}
						`,
						ExpectError: regexp.MustCompile(testCase.ErrorRegex),
					},
				},
			})
		})
	}
}

func TestMailboxResource_Configuration_Success(t *testing.T) {
	tests := []struct {
		testcase      string
		configuration string
		want          model.Mailbox
	}{
		{
			testcase: "invitation",
			configuration: `
				name                    = "Some Name"
				domain_name             = "example.com"
				local_part              = "test"
				password_recovery_email = "someone@example.com"
				password_method         = "invitation"
			`,
			want: model.Mailbox{
				Name:                  "Some Name",
				DomainName:            "example.com",
				LocalPart:             "test",
				Address:               "test@example.com",
				PasswordRecoveryEmail: "someone@example.com",
				PasswordMethod:        "invitation",
			},
		},
		{
			testcase: "managed-password",
			configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
				password    = "secret"
			`,
			want: model.Mailbox{
				Name:       "Some Name",
				DomainName: "example.com",
				LocalPart:  "test",
				Address:    "test@example.com",
				Password:   "secret",
			},
		},
		{
			testcase: "passwords",
			configuration: `
				name                    = "Some Name"
				domain_name             = "example.com"
				local_part              = "test"
				password                = "secret"
				password_recovery_email = "someone@example.com"
			`,
			want: model.Mailbox{
				Name:                  "Some Name",
				DomainName:            "example.com",
				LocalPart:             "test",
				Address:               "test@example.com",
				Password:              "secret",
				PasswordRecoveryEmail: "someone@example.com",
			},
		},
		{
			testcase: "delegations",
			configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
				password    = "secret"
				delegations = ["other@example.com"]
			`,
			want: model.Mailbox{
				Name:        "Some Name",
				DomainName:  "example.com",
				LocalPart:   "test",
				Address:     "test@example.com",
				Password:    "secret",
				Delegations: []string{"other@example.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								%s
							}
						`, tt.configuration),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", tt.want.LocalPart, tt.want.DomainName)),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", tt.want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", tt.want.Password),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password_recovery_email", tt.want.PasswordRecoveryEmail),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "delegations.#", strconv.Itoa(len(tt.want.Delegations))),
						),
					},
				},
			})
		})
	}
}

func TestMailboxResource_Configuration_Errors(t *testing.T) {
	testCases := map[string]ConfigurationErrorTestCase{
		"empty-domain-name": {
			Configuration: `
				name        = "Some Name"
				domain_name = ""
				local_part  = "test"
				password    = "secret"
			`,
			ErrorRegex: "Attribute domain_name string length must be at least 1",
		},
		"missing-domain-name": {
			Configuration: `
				name        = "Some Name"
				local_part  = "test"
				password    = "secret"
			`,
			ErrorRegex: `The argument "domain_name" is required, but no definition was found`,
		},
		"invalid-domain-name": {
			Configuration: `
				name        = "Some Name"
				domain_name = "*.example.com"
				local_part  = "test"
				password    = "secret"
			`,
			ErrorRegex: "Domain names must be convertible to ASCII",
		},
		"empty-local-part": {
			Configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = ""
				password    = "secret"
			`,
			ErrorRegex: "Attribute local_part string length must be at least 1",
		},
		"missing-local-part": {
			Configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				password    = "secret"
			`,
			ErrorRegex: `The argument "local_part" is required, but no definition was found`,
		},
		"empty-password": {
			Configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
				password    = ""
			`,
			ErrorRegex: "Attribute password string length must be at least 1",
		},
		"missing-password": {
			Configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
			`,
			ErrorRegex: `At least one attribute out of \[password\] must be specified`,
		},
		"empty-password-recovery-email": {
			Configuration: `
				name                    = "Some Name"
				domain_name             = "example.com"
				local_part              = "test"
				password_recovery_email = ""
			`,
			ErrorRegex: "Attribute password_recovery_email string length must be at least 1",
		},
		"missing-password-recovery-email": {
			Configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
			`,
			ErrorRegex: `At least one attribute out of \[password_recovery_email\] must be specified`,
		},
		"empty-name": {
			Configuration: `
				name        = ""
				domain_name = "example.com"
				local_part  = "test"
				password    = "secret"
			`,
			ErrorRegex: "Attribute name string length must be at least 1",
		},
		"missing-name": {
			Configuration: `
				domain_name = "example.com"
				local_part  = "test"
				password    = "secret"
			`,
			ErrorRegex: `The argument "name" is required, but no definition was found`,
		},
		"missing-recovery-email": {
			Configuration: `
				name            = "Some Name"
				domain_name     = "example.com"
				local_part      = "test"
				password_method = "invitation"
				password        = "abc"
			`,
			ErrorRegex: "Cannot use 'password_method = invitation' without a 'password_recovery_email'",
		},
		"missing-password-with-password-method": {
			Configuration: `
				name                    = "Some Name"
				domain_name             = "example.com"
				local_part              = "test"
				password_method         = "password"
				password_recovery_email = "someone@example.com"
			`,
			ErrorRegex: "Cannot use 'password_method = password' without a 'password'",
		},
		"unsupported-password-method": {
			Configuration: `
				name            = "Some Name"
				domain_name     = "example.com"
				local_part      = "test"
				password_method = "something"
				password        = "abc"
			`,
			ErrorRegex: "Attribute password_method value must be one of",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								%s
							}
						`, testCase.Configuration),
						ExpectError: regexp.MustCompile(testCase.ErrorRegex),
					},
				},
			})
		})
	}
}
