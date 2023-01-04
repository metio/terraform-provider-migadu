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

func TestAliasesDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		state  []model.Alias
		want   *model.Aliases
	}{
		{
			name:   "single",
			domain: "example.com",
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
			want: &model.Aliases{
				Aliases: []model.Alias{
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
			},
		},
		{
			name:   "multiple",
			domain: "example.com",
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
					Address:          "some@example.com",
					Destinations:     []string{"another@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Aliases{
				Aliases: []model.Alias{
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
						Address:          "some@example.com",
						Destinations:     []string{"another@example"},
						IsInternal:       true,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: false,
					},
				},
			},
		},
		{
			name:   "filtered",
			domain: "example.com",
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
					DomainName:       "different.com",
					Address:          "some@example.com",
					Destinations:     []string{"another@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Aliases{
				Aliases: []model.Alias{
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
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "xn--ho-hia.de",
					Address:          "some@xn--ho-hia.de",
					Destinations:     []string{"other@xn--ho-hia.de"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Aliases{
				Aliases: []model.Alias{
					{
						LocalPart:        "some",
						DomainName:       "xn--ho-hia.de",
						Address:          "some@xn--ho-hia.de",
						Destinations:     []string{"other@xn--ho-hia.de"},
						IsInternal:       true,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: false,
					},
				},
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
							data "migadu_aliases" "test" {
								domain_name = "%s"
							}
						`, tt.domain),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.#", fmt.Sprintf("%v", len(tt.want.Aliases))),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "id", tt.domain),
						),
					},
				},
			})
		})
	}
}

func TestAliasesDataSource_API_Errors(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		error      string
	}{
		{
			name:       "error-404",
			domain:     "example.com",
			statusCode: http.StatusNotFound,
			error:      "GetAliases: status: 404",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			statusCode: http.StatusInternalServerError,
			error:      "GetAliases: status: 500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{StatusCode: tt.statusCode}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							data "migadu_aliases" "test" {
								domain_name = "%s"
							}
						`, tt.domain),
						ExpectError: regexp.MustCompile(tt.error),
					},
				},
			})
		})
	}
}

func TestAliasesDataSource_Configuration_Errors(t *testing.T) {
	tests := []struct {
		name          string
		configuration string
		error         string
	}{
		{
			name: "empty-domain-name",
			configuration: `
				domain_name = ""
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name:          "missing-domain-name",
			configuration: ``,
			error:         `The argument "domain_name" is required, but no definition was found`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							data "migadu_aliases" "test" {
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
