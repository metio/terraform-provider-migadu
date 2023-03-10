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
	"testing"
)

func TestIdentityResource_API_Success(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		localPart   string
		state       []model.Identity
		send        *model.Identity
		updatedName string
		want        *model.Identity
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Identity{},
			send: &model.Identity{
				LocalPart: "other",
				Name:      "Some Name",
				Password:  "secret",
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Some Name",
				Password:   "secret",
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
			send: &model.Identity{
				LocalPart: "other",
				Name:      "Some Name",
				Password:  "secret",
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Some Name",
				Password:   "secret",
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
								domain_name = "%s"
								local_part  = "%s"
								identity    = "%s"
								password    = "%s"
								name        = "%s"
							}
						`, tt.domain, tt.localPart, tt.send.LocalPart, tt.send.Password, tt.send.Name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "identity", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "password", tt.want.Password),
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
							"password", // Migadu API does not allow reading passwords
						},
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_identity" "test" {
								domain_name = "%s"
								local_part  = "%s"
								identity    = "%s"
								password    = "%s"
								name        = "%s"
							}
						`, tt.domain, tt.localPart, tt.send.LocalPart, tt.send.Password, tt.updatedName),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "identity", tt.want.LocalPart),
							resource.TestCheckResourceAttr("migadu_identity.test", "password", tt.want.Password),
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
	tests := []struct {
		name       string
		domain     string
		localPart  string
		identity   string
		password   string
		statusCode int
		state      []model.Identity
		error      string
	}{
		{
			name:      "error-400",
			domain:    "example.com",
			localPart: "test",
			identity:  "someone",
			password:  "secret",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "example.com",
					Address:    "some@example.com",
				},
			},
			error: "CreateIdentity: status:\n400",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			identity:   "identity",
			password:   "secret",
			statusCode: http.StatusInternalServerError,
			error:      "CreateIdentity: status:\n500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_identity" "test" {
								domain_name = "%s"
								local_part  = "%s"
								identity    = "%s"
								password    = "%s"
							}
						`, tt.domain, tt.localPart, tt.identity, tt.password),
						ExpectError: regexp.MustCompile(tt.error),
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
				password    = "secret"
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name: "empty-local-part",
			configuration: `
				domain_name = "example.com"
				local_part  = ""
				identity    = "some"
				password    = "secret"
			`,
			error: "Attribute local_part string length must be at least 1",
		},
		{
			name: "empty-identity",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				identity    = ""
				password    = "secret"
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
			`,
			error: "Attribute password string length must be at least 1",
		},
		{
			name: "missing-domain-name",
			configuration: `
				local_part = "test"
				identity   = "some"
				password    = "secret"
			`,
			error: `The argument "domain_name" is required, but no definition was found`,
		},
		{
			name: "missing-local-part",
			configuration: `
				domain_name = "example.com"
				identity    = "some"
				password    = "secret"
			`,
			error: `The argument "local_part" is required, but no definition was found`,
		},
		{
			name: "missing-identity",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				password    = "secret"
			`,
			error: `The argument "identity" is required, but no definition was found`,
		},
		{
			name: "missing-password",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				identity    = "some"
			`,
			error: `The argument "password" is required, but no definition was found`,
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
