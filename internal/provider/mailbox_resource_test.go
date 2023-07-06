//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"github.com/metio/terraform-provider-migadu/migadu/simulator"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
)

func TestMailboxResource_API_Success_Using_Password(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		state       []model.Mailbox
		send        *model.Mailbox
		updatedName string
		want        *model.Mailbox
	}{
		{
			name:   "single",
			domain: "example.com",
			state:  []model.Mailbox{},
			send: &model.Mailbox{
				LocalPart: "test",
				Name:      "Some Name",
				Password:  "secret",
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
			updatedName: "Different Name",
		},
		{
			name:   "multiple",
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "different.com",
					Address:    "test@different.com",
					Name:       "Some Name",
				},
			},
			send: &model.Mailbox{
				LocalPart: "test",
				Name:      "Some Name",
				Password:  "secret",
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
			updatedName: "Different Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								domain_name = "%s"
								local_part  = "%s"
								password    = "%s"
								name        = "%s"
							}
						`, tt.domain, tt.send.LocalPart, tt.send.Password, tt.send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", tt.want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", tt.send.Password),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", tt.want.LocalPart, tt.want.DomainName)),
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
								domain_name = "%s"
								local_part  = "%s"
								password    = "%s"
								name        = "%s"
							}
						`, tt.domain, tt.send.LocalPart, tt.send.Password, tt.updatedName),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", tt.want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", tt.updatedName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", tt.send.Password),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", tt.want.LocalPart, tt.want.DomainName)),
						),
					},
				},
			})
		})
	}
}

func TestMailboxResource_API_Success_Using_RecoveryEmail(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		state       []model.Mailbox
		send        *model.Mailbox
		updatedName string
		want        *model.Mailbox
	}{
		{
			name:   "single",
			domain: "example.com",
			state:  []model.Mailbox{},
			send: &model.Mailbox{
				LocalPart:             "test",
				Name:                  "Some Name",
				PasswordRecoveryEmail: "someone@example.com",
			},
			want: &model.Mailbox{
				LocalPart:             "test",
				DomainName:            "example.com",
				Address:               "test@example.com",
				Name:                  "Some Name",
				PasswordRecoveryEmail: "someone@example.com",
			},
			updatedName: "Different Name",
		},
		{
			name:   "multiple",
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "different.com",
					Address:    "test@different.com",
					Name:       "Some Name",
				},
			},
			send: &model.Mailbox{
				LocalPart:             "test",
				Name:                  "Some Name",
				PasswordRecoveryEmail: "someone@example.com",
			},
			want: &model.Mailbox{
				LocalPart:             "test",
				DomainName:            "example.com",
				Address:               "test@example.com",
				Name:                  "Some Name",
				PasswordRecoveryEmail: "someone@example.com",
			},
			updatedName: "Different Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								domain_name             = "%s"
								local_part              = "%s"
								password_recovery_email = "%s"
								name                    = "%s"
							}
						`, tt.domain, tt.send.LocalPart, tt.send.PasswordRecoveryEmail, tt.send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", tt.want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", ""),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password_recovery_email", tt.want.PasswordRecoveryEmail),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", tt.want.LocalPart, tt.want.DomainName)),
						),
					},
					{
						ResourceName:      "migadu_mailbox.test",
						ImportState:       true,
						ImportStateVerify: true,
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								domain_name             = "%s"
								local_part              = "%s"
								password_recovery_email = "%s"
								name                    = "%s"
							}
						`, tt.domain, tt.send.LocalPart, tt.send.PasswordRecoveryEmail, tt.updatedName),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", tt.want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", tt.updatedName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", ""),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password_recovery_email", tt.want.PasswordRecoveryEmail),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", tt.want.LocalPart, tt.want.DomainName)),
						),
					},
				},
			})
		})
	}
}

func TestMailboxResource_API_Errors(t *testing.T) {
	tests := []struct {
		testcase   string
		name       string
		domain     string
		localPart  string
		password   string
		statusCode int
		state      []model.Mailbox
		error      string
	}{
		{
			testcase:  "error-400",
			name:      "Some Name",
			domain:    "example.com",
			localPart: "test",
			password:  "secret",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
				},
			},
			error: "CreateMailbox: status: 400",
		},
		{
			testcase:   "error-500",
			name:       "Some Name",
			domain:     "example.com",
			localPart:  "test",
			password:   "secret",
			statusCode: http.StatusInternalServerError,
			error:      "CreateMailbox: status: 500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								name        = "%s"
								domain_name = "%s"
								local_part  = "%s"
								password    = "%s"
							}
						`, tt.name, tt.domain, tt.localPart, tt.password),
						ExpectError: regexp.MustCompile(tt.error),
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
		state         []model.Mailbox
		want          *model.Mailbox
	}{
		{
			testcase: "invitation",
			configuration: `
				name                    = "Some Name"
				domain_name             = "example.com"
				local_part              = "test"
				password_recovery_email = "someone@example.com"
			`,
			state: []model.Mailbox{},
			want: &model.Mailbox{
				Name:                  "Some Name",
				DomainName:            "example.com",
				LocalPart:             "test",
				Address:               "test@example.com",
				PasswordRecoveryEmail: "someone@example.com",
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
			state: []model.Mailbox{},
			want: &model.Mailbox{
				Name:       "Some Name",
				DomainName: "example.com",
				LocalPart:  "test",
				Address:    "test@example.com",
				Password:   "secret",
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
			state: []model.Mailbox{},
			want: &model.Mailbox{
				Name:        "Some Name",
				DomainName:  "example.com",
				LocalPart:   "test",
				Address:     "test@example.com",
				Password:    "secret",
				Delegations: []string{"other@example.com"},
			},
		},
		{
			testcase: "identities",
			configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
				password    = "secret"
				identities = ["other@example.com"]
			`,
			state: []model.Mailbox{},
			want: &model.Mailbox{
				Name:       "Some Name",
				DomainName: "example.com",
				LocalPart:  "test",
				Address:    "test@example.com",
				Password:   "secret",
				Identities: []string{"other@example.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state}))
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
							resource.TestCheckResourceAttr("migadu_mailbox.test", "domain_name", tt.want.DomainName),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "local_part", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password", tt.want.Password),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "password_recovery_email", tt.want.PasswordRecoveryEmail),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", tt.want.LocalPart, tt.want.DomainName)),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "delegations.#", strconv.Itoa(len(tt.want.Delegations))),
							resource.TestCheckResourceAttr("migadu_mailbox.test", "identities.#", strconv.Itoa(len(tt.want.Identities))),
						),
					},
				},
			})
		})
	}
}

func TestMailboxResource_Configuration_Errors(t *testing.T) {
	tests := []struct {
		name          string
		configuration string
		error         string
	}{
		{
			name: "empty-domain-name",
			configuration: `
				name        = "Some Name"
				domain_name = ""
				local_part  = "test"
				password    = "secret"
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name: "missing-domain-name",
			configuration: `
				name        = "Some Name"
				local_part  = "test"
				password    = "secret"
			`,
			error: `The argument "domain_name" is required, but no definition was found`,
		},
		{
			name: "empty-local-part",
			configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = ""
				password    = "secret"
			`,
			error: "Attribute local_part string length must be at least 1",
		},
		{
			name: "missing-local-part",
			configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				password    = "secret"
			`,
			error: `The argument "local_part" is required, but no definition was found`,
		},
		{
			name: "empty-password",
			configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
				password    = ""
			`,
			error: "Attribute password string length must be at least 1",
		},
		{
			name: "missing-password",
			configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
			`,
			error: `No attribute specified when one \(and only one\) of \[password\] is required`,
		},
		{
			name: "empty-password-recovery-email",
			configuration: `
				name                    = "Some Name"
				domain_name             = "example.com"
				local_part              = "test"
				password_recovery_email = ""
			`,
			error: "Attribute password_recovery_email string length must be at least 1",
		},
		{
			name: "missing-password-recovery-email",
			configuration: `
				name        = "Some Name"
				domain_name = "example.com"
				local_part  = "test"
			`,
			error: `(?s)No attribute specified when one \(and only one\) of \[password_recovery_email\](.*)is required`,
		},
		{
			name: "empty-name",
			configuration: `
				name        = ""
				domain_name = "example.com"
				local_part  = "test"
				password    = "secret"
			`,
			error: "Attribute name string length must be at least 1",
		},
		{
			name: "missing-name",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				password    = "secret"
			`,
			error: `The argument "name" is required, but no definition was found`,
		},
		{
			name: "duplicate-passwords",
			configuration: `
				name                    = "Some Name"
				domain_name             = "example.com"
				local_part              = "test"
				password                = "secret"
				password_recovery_email = "someone@example.com"
			`,
			error: `2 attributes specified when one \(and only one\) of \[password\] is required`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							resource "migadu_mailbox" "test" {
								%s
							}
						`, tt.configuration),
						ExpectError: regexp.MustCompile(tt.error),
					},
				},
			})
		})
	}
}
