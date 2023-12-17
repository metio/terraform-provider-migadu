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
	"github.com/metio/migadu-client.go/model"
	"github.com/metio/migadu-client.go/simulator"
	"github.com/metio/terraform-provider-migadu/internal/provider"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestRewriteRulesDataSource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	provider.NewRewriteRulesDataSource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestRewriteRulesDataSource_API_Success(t *testing.T) {
	testCases := map[string]struct {
		domain string
		state  []model.RewriteRule
		want   model.RewriteRules
	}{
		"single": {
			domain: "example.com",
			state: []model.RewriteRule{
				{
					DomainName:    "example.com",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@example.com",
					},
				},
			},
			want: model.RewriteRules{
				RewriteRules: []model.RewriteRule{
					{
						DomainName:    "example.com",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest@example.com",
						},
					},
				},
			},
		},
		"multiple": {
			domain: "example.com",
			state: []model.RewriteRule{
				{
					DomainName:    "example.com",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest1@example.com",
						"dest2@example.com",
					},
				},
				{
					DomainName:    "example.com",
					Name:          "Other Rule",
					LocalPartRule: "*-suffix",
					OrderNum:      1,
					Destinations: []string{
						"dest3@example.com",
						"dest4@example.com",
					},
				},
			},
			want: model.RewriteRules{
				RewriteRules: []model.RewriteRule{
					{
						DomainName:    "example.com",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest1@example.com",
							"dest2@example.com",
						},
					},
					{
						DomainName:    "example.com",
						Name:          "Other Rule",
						LocalPartRule: "*-suffix",
						OrderNum:      1,
						Destinations: []string{
							"dest3@example.com",
							"dest4@example.com",
						},
					},
				},
			},
		},
		"filtered": {
			domain: "example.com",
			state: []model.RewriteRule{
				{
					DomainName:    "different.com",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest1@different.com",
						"dest2@different.com",
					},
				},
				{
					DomainName:    "example.com",
					Name:          "Other Rule",
					LocalPartRule: "*-suffix",
					OrderNum:      1,
					Destinations: []string{
						"dest3@example.com",
						"dest4@example.com",
					},
				},
			},
			want: model.RewriteRules{
				RewriteRules: []model.RewriteRule{
					{
						DomainName:    "example.com",
						Name:          "Other Rule",
						LocalPartRule: "*-suffix",
						OrderNum:      1,
						Destinations: []string{
							"dest3@example.com",
							"dest4@example.com",
						},
					},
				},
			},
		},
		"idna": {
			domain: "hoß.de",
			state: []model.RewriteRule{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@xn--ho-hia.de",
					},
				},
			},
			want: model.RewriteRules{
				RewriteRules: []model.RewriteRule{
					{
						DomainName:    "xn--ho-hia.de",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest@xn--ho-hia.de",
						},
					},
				},
			},
		},
		"idna-punycode": {
			domain: "xn--ho-hia.de",
			state: []model.RewriteRule{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@xn--ho-hia.de",
					},
				},
			},
			want: model.RewriteRules{
				RewriteRules: []model.RewriteRule{
					{
						DomainName:    "xn--ho-hia.de",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest@xn--ho-hia.de",
						},
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: testCase.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							data "migadu_rewrite_rules" "test" {
								domain_name = "%s"
							}
						`, testCase.domain),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_rewrite_rules.test", "id", testCase.domain),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rules.test", "domain_name", testCase.domain),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rules.test", "rewrites.#", fmt.Sprintf("%v", len(testCase.want.RewriteRules))),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rules.test", "rewrites.0.local_part_rule", testCase.want.RewriteRules[0].LocalPartRule),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rules.test", "rewrites.0.order_num", fmt.Sprintf("%v", testCase.want.RewriteRules[0].OrderNum)),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rules.test", "rewrites.0.destinations.#", fmt.Sprintf("%v", len(testCase.want.RewriteRules[0].Destinations))),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rules.test", "rewrites.0.destinations.0", testCase.want.RewriteRules[0].Destinations[0]),
						),
					},
				},
			})
		})
	}
}

func TestRewriteRulesDataSource_API_Errors(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "GetRewriteRules: status: 404",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "GetRewriteRules: status: 500",
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
							data "migadu_rewrite_rules" "test" {
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

func TestRewriteRulesDataSource_Configuration_Errors(t *testing.T) {
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
							data "migadu_rewrite_rules" "test" {
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
