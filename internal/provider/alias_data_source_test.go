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

func TestAliasDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name      string
		domain    string
		localPart string
		state     []model.Alias
		want      *model.Alias
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "some",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Destinations:     []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:        "some",
				DomainName:       "example.com",
				Address:          "some@example.com",
				Destinations:     []string{"other@example"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "some",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Destinations:     []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
				{
					LocalPart:        "other",
					DomainName:       "example.com",
					Address:          "other@example.com",
					Destinations:     []string{"different@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:        "some",
				DomainName:       "example.com",
				Address:          "some@example.com",
				Destinations:     []string{"other@example"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:        "test",
					DomainName:       "xn--ho-hia.de",
					Address:          "test@xn--ho-hia.de",
					Destinations:     []string{"other@xn--ho-hia.de"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:        "test",
				DomainName:       "xn--ho-hia.de",
				Address:          "test@xn--ho-hia.de",
				Destinations:     []string{"other@xn--ho-hia.de"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Aliases: tt.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							data "migadu_alias" "test" {
								domain_name = "%s"
								local_part  = "%s"
							}
						`, tt.domain, tt.localPart),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_alias.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "address", tt.want.Address),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "destinations.#", fmt.Sprintf("%v", len(tt.want.Destinations))),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "destinations_punycode.#", fmt.Sprintf("%v", len(tt.want.Destinations))),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "id", fmt.Sprintf("%s@%s", tt.localPart, tt.domain)),
						),
					},
				},
			})
		})
	}
}

func TestAliasDataSource_API_Errors(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Alias
		error      string
	}{
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "other",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Destinations:     []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			error: "GetAlias: status: 404",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "other",
			statusCode: http.StatusInternalServerError,
			error:      "GetAlias: status: 500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Aliases: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							data "migadu_alias" "test" {
								domain_name = "%s"
								local_part  = "%s"
							}
						`, tt.domain, tt.localPart),
						ExpectError: regexp.MustCompile(tt.error),
					},
				},
			})
		})
	}
}

func TestAliasDataSource_Configuration_Errors(t *testing.T) {
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
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name: "empty-local-part",
			configuration: `
				domain_name = "example.com"
				local_part  = ""
			`,
			error: "Attribute local_part string length must be at least 1",
		},
		{
			name: "missing-domain-name",
			configuration: `
				local_part  = "test"
			`,
			error: `The argument "domain_name" is required, but no definition was found`,
		},
		{
			name: "missing-local-part",
			configuration: `
				domain_name = "example.com"
			`,
			error: `The argument "local_part" is required, but no definition was found`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							data "migadu_alias" "test" {
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
