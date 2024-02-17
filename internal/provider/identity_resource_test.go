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
	"testing"
)

func TestIdentityResource_API_Success_With_Password(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		localPart   string
		state       []model.Identity
		send        model.Identity
		updatedName string
		want        model.Identity
	}{
		{
			name:      "single-custom",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Identity{},
			send: model.Identity{
				LocalPart:   "other",
				Name:        "Some Name",
				Password:    "secret",
				PasswordUse: "custom",
			},
			want: model.Identity{
				LocalPart:   "other",
				DomainName:  "example.com",
				Address:     "other@example.com",
				Name:        "Some Name",
				Password:    "secret",
				PasswordUse: "custom",
			},
			updatedName: "Different Name",
		},
		{
			name:      "multiple-custom",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "example.com",
					Address:    "some@example.com",
				},
			},
			send: model.Identity{
				LocalPart:   "other",
				Name:        "Some Name",
				Password:    "secret",
				PasswordUse: "custom",
			},
			want: model.Identity{
				LocalPart:   "other",
				DomainName:  "example.com",
				Address:     "other@example.com",
				Name:        "Some Name",
				Password:    "secret",
				PasswordUse: "custom",
			},
			updatedName: "Different Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_identity" "test" {
								domain_name  = "%s"
								local_part   = "%s"
								identity     = "%s"
								password     = "%s"
								password_use = "%s"
								name         = "%s"
							}
						`, tt.domain, tt.localPart, tt.send.LocalPart, tt.send.Password, tt.send.PasswordUse, tt.send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "identity", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "password", tt.want.Password),
							resource.TestCheckResourceAttr("migadu_identity.test", "password_use", tt.want.PasswordUse),
							resource.TestCheckResourceAttr("migadu_identity.test", "address", tt.want.Address),
							resource.TestCheckResourceAttr("migadu_identity.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_identity.test", "id", fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.send.LocalPart)),
						),
					},
					{
						ResourceName:      "migadu_identity.test",
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateVerifyIgnore: []string{
							"password",     // Migadu API does not allow reading passwords
							"password_use", // Migadu API does not allow reading passwords
						},
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_identity" "test" {
								domain_name  = "%s"
								local_part   = "%s"
								identity     = "%s"
								password     = "%s"
								password_use = "%s"
								name         = "%s"
							}
						`, tt.domain, tt.localPart, tt.send.LocalPart, tt.send.Password, tt.send.PasswordUse, tt.updatedName),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "identity", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "password", tt.want.Password),
							resource.TestCheckResourceAttr("migadu_identity.test", "password_use", tt.want.PasswordUse),
							resource.TestCheckResourceAttr("migadu_identity.test", "address", tt.want.Address),
							resource.TestCheckResourceAttr("migadu_identity.test", "name", tt.updatedName),
							resource.TestCheckResourceAttr("migadu_identity.test", "id", fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.send.LocalPart)),
						),
					},
				},
			})
		})
	}
}

func TestIdentityResource_API_Success_Without_Password(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		localPart   string
		state       []model.Identity
		send        model.Identity
		updatedName string
		want        model.Identity
	}{
		{
			name:      "single-none",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Identity{},
			send: model.Identity{
				LocalPart:   "other",
				Name:        "Some Name",
				PasswordUse: "none",
			},
			want: model.Identity{
				LocalPart:   "other",
				DomainName:  "example.com",
				Address:     "other@example.com",
				Name:        "Some Name",
				PasswordUse: "none",
			},
			updatedName: "Different Name",
		},
		{
			name:      "multiple-none",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "example.com",
					Address:    "some@example.com",
				},
			},
			send: model.Identity{
				LocalPart:   "other",
				Name:        "Some Name",
				PasswordUse: "none",
			},
			want: model.Identity{
				LocalPart:   "other",
				DomainName:  "example.com",
				Address:     "other@example.com",
				Name:        "Some Name",
				PasswordUse: "none",
			},
			updatedName: "Different Name",
		},
		{
			name:      "single-mailbox",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Identity{},
			send: model.Identity{
				LocalPart:   "other",
				Name:        "Some Name",
				PasswordUse: "mailbox",
			},
			want: model.Identity{
				LocalPart:   "other",
				DomainName:  "example.com",
				Address:     "other@example.com",
				Name:        "Some Name",
				PasswordUse: "mailbox",
			},
			updatedName: "Different Name",
		},
		{
			name:      "multiple-mailbox",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "example.com",
					Address:    "some@example.com",
				},
			},
			send: model.Identity{
				LocalPart:   "other",
				Name:        "Some Name",
				PasswordUse: "mailbox",
			},
			want: model.Identity{
				LocalPart:   "other",
				DomainName:  "example.com",
				Address:     "other@example.com",
				Name:        "Some Name",
				PasswordUse: "mailbox",
			},
			updatedName: "Different Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_identity" "test" {
								domain_name  = "%s"
								local_part   = "%s"
								identity     = "%s"
								password_use = "%s"
								name         = "%s"
							}
						`, tt.domain, tt.localPart, tt.send.LocalPart, tt.send.PasswordUse, tt.send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "identity", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "password_use", tt.want.PasswordUse),
							resource.TestCheckResourceAttr("migadu_identity.test", "address", tt.want.Address),
							resource.TestCheckResourceAttr("migadu_identity.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_identity.test", "id", fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.send.LocalPart)),
						),
					},
					{
						ResourceName:      "migadu_identity.test",
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateVerifyIgnore: []string{
							"password_use", // Migadu API does not allow reading passwords
						},
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_identity" "test" {
								domain_name  = "%s"
								local_part   = "%s"
								identity     = "%s"
								password_use = "%s"
								name         = "%s"
							}
						`, tt.domain, tt.localPart, tt.send.LocalPart, tt.send.PasswordUse, tt.updatedName),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "identity", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "password_use", tt.want.PasswordUse),
							resource.TestCheckResourceAttr("migadu_identity.test", "address", tt.want.Address),
							resource.TestCheckResourceAttr("migadu_identity.test", "name", tt.updatedName),
							resource.TestCheckResourceAttr("migadu_identity.test", "id", fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.send.LocalPart)),
						),
					},
				},
			})
		})
	}
}

func TestIdentityResource_API_Success_With_Default_PasswordUse(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		localPart   string
		state       []model.Identity
		send        model.Identity
		updatedName string
		want        model.Identity
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Identity{},
			send: model.Identity{
				LocalPart:   "other",
				Name:        "Some Name",
				PasswordUse: "none",
			},
			want: model.Identity{
				LocalPart:   "other",
				DomainName:  "example.com",
				Address:     "other@example.com",
				Name:        "Some Name",
				PasswordUse: "none",
			},
			updatedName: "Different Name",
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "example.com",
					Address:    "some@example.com",
				},
			},
			send: model.Identity{
				LocalPart:   "other",
				Name:        "Some Name",
				PasswordUse: "none",
			},
			want: model.Identity{
				LocalPart:   "other",
				DomainName:  "example.com",
				Address:     "other@example.com",
				Name:        "Some Name",
				PasswordUse: "none",
			},
			updatedName: "Different Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_identity" "test" {
								domain_name  = "%s"
								local_part   = "%s"
								identity     = "%s"
								name         = "%s"
							}
						`, tt.domain, tt.localPart, tt.send.LocalPart, tt.send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "identity", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "password_use", tt.want.PasswordUse),
							resource.TestCheckResourceAttr("migadu_identity.test", "address", tt.want.Address),
							resource.TestCheckResourceAttr("migadu_identity.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_identity.test", "id", fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.send.LocalPart)),
						),
					},
					{
						ResourceName:      "migadu_identity.test",
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateVerifyIgnore: []string{
							"password_use", // Migadu API does not allow reading passwords
						},
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_identity" "test" {
								domain_name  = "%s"
								local_part   = "%s"
								identity     = "%s"
								name         = "%s"
							}
						`, tt.domain, tt.localPart, tt.send.LocalPart, tt.updatedName),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "identity", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "password_use", tt.want.PasswordUse),
							resource.TestCheckResourceAttr("migadu_identity.test", "address", tt.want.Address),
							resource.TestCheckResourceAttr("migadu_identity.test", "name", tt.updatedName),
							resource.TestCheckResourceAttr("migadu_identity.test", "id", fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.send.LocalPart)),
						),
					},
				},
			})
		})
	}
}

func TestIdentityResource_API_Errors(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "CreateIdentity: status: 404",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "CreateIdentity: status: 500",
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
							resource "migadu_identity" "test" {
								local_part   = "test"
								domain_name  = "example.com"
								identity     = "someone"
								name         = "Some Name"
								password     = "supers3cret"
								password_use = "custom"
							}
						`,
						ExpectError: regexp.MustCompile(testCase.ErrorRegex),
					},
				},
			})
		})
	}
}

func TestIdentityResource_Configuration_Errors(t *testing.T) {
	tests := []struct {
		name          string
		configuration string
		error         string
	}{
		{
			name: "empty-domain-name",
			configuration: `
				domain_name = ""
				local_part  = "test"
				identity    = "some"
				name        = "Some Name"
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name: "empty-local-part",
			configuration: `
				domain_name = "example.com"
				local_part  = ""
				identity    = "some"
				name        = "Some Name"
			`,
			error: "Attribute local_part string length must be at least 1",
		},
		{
			name: "empty-identity",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				identity    = ""
				name        = "Some Name"
			`,
			error: "Attribute identity string length must be at least 1",
		},
		{
			name: "empty-password",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				identity    = "some"
				password    = ""
				name        = "Some Name"
			`,
			error: "Attribute password string length must be at least 1",
		},
		{
			name: "missing-domain-name",
			configuration: `
				local_part = "test"
				identity   = "some"
				name        = "Some Name"
			`,
			error: `The argument "domain_name" is required, but no definition was found`,
		},
		{
			name: "missing-local-part",
			configuration: `
				domain_name = "example.com"
				identity    = "some"
				name        = "Some Name"
			`,
			error: `The argument "local_part" is required, but no definition was found`,
		},
		{
			name: "missing-identity",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				name        = "Some Name"
			`,
			error: `The argument "identity" is required, but no definition was found`,
		},
		{
			name: "missing-name",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				identity    = "some"
			`,
			error: `The argument "name" is required, but no definition was found`,
		},
		{
			name: "empty-name",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				identity    = "some"
				name        = ""
			`,
			error: `Attribute name string length must be at least 1, got: 0`,
		},
		{
			name: "wrong-password-use",
			configuration: `
				domain_name  = "example.com"
				local_part   = "test"
				identity     = "some"
				name         = "Some Name"
				password     = "secret"
				password_use = "invalid"
			`,
			error: `Attribute password_use value must be one of: \["none" "mailbox" "custom"\]`,
		},
		{
			name: "no-custom-password",
			configuration: `
				domain_name  = "example.com"
				local_part   = "test"
				identity     = "some"
				name         = "Some Name"
				password_use = "custom"
			`,
			error: `Attribute "password" must be specified when "password_use" is specified`,
		},
		{
			name: "unnecessary-none-password",
			configuration: `
				domain_name  = "example.com"
				local_part   = "test"
				identity     = "some"
				name         = "Some Name"
				password_use = "none"
				password     = "secret"
			`,
			error: `Attribute "password" cannot be specified when "password_use" is specified`,
		},
		{
			name: "unnecessary-mailbox-password",
			configuration: `
				domain_name  = "example.com"
				local_part   = "test"
				identity     = "some"
				name         = "Some Name"
				password_use = "mailbox"
				password     = "secret"
			`,
			error: `Attribute "password" cannot be specified when "password_use" is specified`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							resource "migadu_identity" "test" {
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
