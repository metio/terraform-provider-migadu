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

func TestIdentityDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name      string
		domain    string
		localPart string
		identity  string
		state     []model.Identity
		want      *model.Identity
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			identity:  "someone",
			state: []model.Identity{
				{
					LocalPart:            "someone",
					DomainName:           "example.com",
					Address:              "someone@example.com",
					Name:                 "Some Identity",
					MaySend:              true,
					MayReceive:           true,
					MayAccessImap:        true,
					MayAccessPop3:        true,
					MayAccessManageSieve: true,
					Password:             "secret",
					FooterActive:         false,
					FooterPlainBody:      "",
					FooterHtmlBody:       "",
				},
			},
			want: &model.Identity{
				LocalPart:            "someone",
				DomainName:           "example.com",
				Address:              "someone@example.com",
				Name:                 "Some Identity",
				MaySend:              true,
				MayReceive:           true,
				MayAccessImap:        true,
				MayAccessPop3:        true,
				MayAccessManageSieve: true,
				Password:             "secret",
				FooterActive:         false,
				FooterPlainBody:      "",
				FooterHtmlBody:       "",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			identity:  "someone",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "example.com",
					Address:    "someone@example.com",
					Name:       "Some Identity",
				},
				{
					LocalPart:  "another",
					DomainName: "example.com",
					Address:    "another@example.com",
					Name:       "Another Identity",
				},
			},
			want: &model.Identity{
				LocalPart:  "someone",
				DomainName: "example.com",
				Address:    "someone@example.com",
				Name:       "Some Identity",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			identity:  "someone",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "xn--ho-hia.de",
					Address:    "someone@xn--ho-hia.de",
					Name:       "Some Identity",
				},
			},
			want: &model.Identity{
				LocalPart:  "someone",
				DomainName: "xn--ho-hia.de",
				Address:    "someone@xn--ho-hia.de",
			},
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
							data "migadu_identity" "test" {
								domain_name = "%s"
								local_part  = "%s"
								identity    = "%s"
							}
						`, tt.domain, tt.localPart, tt.identity),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_identity.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("data.migadu_identity.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("data.migadu_identity.test", "identity", tt.identity),
							resource.TestCheckResourceAttr("data.migadu_identity.test", "id", fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.identity)),
						),
					},
				},
			})
		})
	}
}

func TestIdentityDataSource_API_Errors(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		identity   string
		statusCode int
		state      []model.Identity
		error      string
	}{
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			identity:  "someone",
			state: []model.Identity{
				{
					LocalPart:  "some",
					DomainName: "example.com",
					Address:    "some@example.com",
				},
			},
			error: "GetIdentity: status: 404",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			identity:   "identity",
			statusCode: http.StatusInternalServerError,
			error:      "GetIdentity: status: 500",
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
							data "migadu_identity" "test" {
								domain_name = "%s"
								local_part  = "%s"
								identity    = "%s"
							}
						`, tt.domain, tt.localPart, tt.identity),
						ExpectError: regexp.MustCompile(tt.error),
					},
				},
			})
		})
	}
}

func TestIdentityDataSource_Configuration_Errors(t *testing.T) {
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
				identity    = "test"
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name: "empty-local-part",
			configuration: `
				domain_name = "example.com"
				local_part  = ""
				identity    = "test"
			`,
			error: "Attribute local_part string length must be at least 1",
		},
		{
			name: "empty-identity",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
				identity    = ""
			`,
			error: "Attribute identity string length must be at least 1",
		},
		{
			name: "missing-domain-name",
			configuration: `
				local_part  = "test"
				identity    = "test"
			`,
			error: `The argument "domain_name" is required, but no definition was found`,
		},
		{
			name: "missing-local-part",
			configuration: `
				domain_name = "example.com"
				identity    = "test"
			`,
			error: `The argument "local_part" is required, but no definition was found`,
		},
		{
			name: "missing-identity",
			configuration: `
				domain_name = "example.com"
				local_part  = "test"
			`,
			error: `The argument "identity" is required, but no definition was found`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							data "migadu_identity" "test" {
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
