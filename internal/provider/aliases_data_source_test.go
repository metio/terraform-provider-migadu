//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"context"
	"fmt"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-migadu/internal/provider"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"github.com/metio/terraform-provider-migadu/migadu/simulator"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestAliasesDataSource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	provider.NewAliasesDataSource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestAliasesDataSource_API_Success(t *testing.T) {
	testCases := map[string]struct {
		domain string
		state  []model.Alias
		want   model.Aliases
	}{
		"single": {
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
			want: model.Aliases{
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
		"multiple": {
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
			want: model.Aliases{
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
		"filtered": {
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
			want: model.Aliases{
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
		"idna": {
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
			want: model.Aliases{
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
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Aliases: testCase.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							data "migadu_aliases" "test" {
								domain_name = "%s"
							}
						`, testCase.domain),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "id", testCase.domain),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "domain_name", testCase.domain),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.#", fmt.Sprintf("%v", len(testCase.want.Aliases))),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.local_part", testCase.want.Aliases[0].LocalPart),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.domain_name", testCase.want.Aliases[0].DomainName),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.address", testCase.want.Aliases[0].Address),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.destinations.#", fmt.Sprintf("%v", len(testCase.want.Aliases[0].Destinations))),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.destinations.0", testCase.want.Aliases[0].Destinations[0]),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.is_internal", fmt.Sprintf("%v", testCase.want.Aliases[0].IsInternal)),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.expirable", fmt.Sprintf("%v", testCase.want.Aliases[0].Expirable)),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.expires_on", testCase.want.Aliases[0].ExpiresOn),
							resource.TestCheckResourceAttr("data.migadu_aliases.test", "aliases.0.remove_upon_expiry", fmt.Sprintf("%v", testCase.want.Aliases[0].RemoveUponExpiry)),
						),
					},
				},
			})
		})
	}
}

func TestAliasesDataSource_API_Errors(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "GetAliases: status: 404",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "GetAliases: status: 500",
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
							data "migadu_aliases" "test" {
								domain_name = "example.com"
							}
						`,
						ExpectError: regexp.MustCompile(testCase.ErrorRegex),
					},
				},
			})
		})
	}
}

func TestAliasesDataSource_Configuration_Errors(t *testing.T) {
	testCases := map[string]ConfigurationErrorTestCase{
		"empty-domain-name": {
			Configuration: `
				domain_name = ""
			`,
			ErrorRegex: "Attribute domain_name string length must be at least 1",
		},
		"missing-domain-name": {
			Configuration: ``,
			ErrorRegex:    `The argument "domain_name" is required, but no definition was found`,
		},
		"invalid-domain-name": {
			Configuration: `
				domain_name = "*.example.com"
			`,
			ErrorRegex: "Domain names must be convertible to ASCII",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							data "migadu_aliases" "test" {
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
