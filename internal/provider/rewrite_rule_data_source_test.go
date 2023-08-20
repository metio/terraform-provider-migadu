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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-migadu/internal/provider"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"github.com/metio/terraform-provider-migadu/migadu/simulator"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestRewriteRuleDataSource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	provider.NewRewriteRuleDataSource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestRewriteRuleDataSource_API_Success(t *testing.T) {
	testCases := map[string]struct {
		domain string
		name   string
		state  []model.RewriteRule
		want   model.RewriteRule
	}{
		"single": {
			domain: "example.com",
			name:   "test",
			state: []model.RewriteRule{
				{
					DomainName:    "example.com",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      5,
					Destinations: []string{
						"dest@example.com",
					},
				},
			},
			want: model.RewriteRule{
				LocalPartRule: "prefix-*",
				OrderNum:      5,
				Destinations: []string{
					"dest@example.com",
				},
			},
		},
		"multiple": {
			domain: "example.com",
			name:   "test",
			state: []model.RewriteRule{
				{
					DomainName:    "different.com",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      1,
					Destinations: []string{
						"dest@different.com",
					},
				},
				{
					DomainName:    "example.com",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      3,
					Destinations: []string{
						"dest@example.com",
					},
				},
			},
			want: model.RewriteRule{
				LocalPartRule: "prefix-*",
				OrderNum:      3,
				Destinations: []string{
					"dest@example.com",
				},
			},
		},
		"idna": {
			domain: "hoß.de",
			name:   "test",
			state: []model.RewriteRule{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      2,
					Destinations: []string{
						"dest@xn--ho-hia.de",
					},
				},
			},
			want: model.RewriteRule{
				LocalPartRule: "prefix-*",
				OrderNum:      2,
				Destinations: []string{
					"dest@xn--ho-hia.de",
				},
			},
		},
		"idna-punycode": {
			domain: "xn--ho-hia.de",
			name:   "test",
			state: []model.RewriteRule{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      2,
					Destinations: []string{
						"dest@hoß.de",
					},
				},
			},
			want: model.RewriteRule{
				LocalPartRule: "prefix-*",
				OrderNum:      2,
				Destinations: []string{
					"dest@hoß.de",
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
							data "migadu_rewrite_rule" "test" {
								domain_name = "%s"
								name        = "%s"
							}
						`, testCase.domain, testCase.name),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_rewrite_rule.test", "id", fmt.Sprintf("%s/%s", testCase.domain, testCase.name)),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rule.test", "domain_name", testCase.domain),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rule.test", "name", testCase.name),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rule.test", "local_part_rule", testCase.want.LocalPartRule),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rule.test", "order_num", fmt.Sprintf("%v", testCase.want.OrderNum)),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rule.test", "destinations.#", fmt.Sprintf("%v", len(testCase.want.Destinations))),
							resource.TestCheckResourceAttr("data.migadu_rewrite_rule.test", "destinations.0", testCase.want.Destinations[0]),
						),
					},
				},
			})
		})
	}
}

func TestRewriteRuleDataSource_API_Errors(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "GetRewriteRule: status: 404",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "GetRewriteRule: status: 500",
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
							data "migadu_rewrite_rule" "test" {
								domain_name = "example.com"
								name        = "some"
							}
						`,
						ExpectError: regexp.MustCompile(testCase.ErrorRegex),
					},
				},
			})
		})
	}
}

func TestRewriteRuleDataSource_Configuration_Errors(t *testing.T) {
	testCases := map[string]ConfigurationErrorTestCase{
		"empty-domain-name": {
			Configuration: `
				domain_name = ""
				name        = "some-name"
			`,
			ErrorRegex: "Attribute domain_name string length must be at least 1",
		},
		"empty-name": {
			Configuration: `
				domain_name = "example.com"
				name        = ""
			`,
			ErrorRegex: "Attribute name string length must be at least 1",
		},
		"missing-domain-name": {
			Configuration: `
				name = "some-name"
			`,
			ErrorRegex: `The argument "domain_name" is required, but no definition was found`,
		},
		"missing-name": {
			Configuration: `
				domain_name = "example.com"
			`,
			ErrorRegex: `The argument "name" is required, but no definition was found`,
		},
		"invalid-domain-name": {
			Configuration: `
				domain_name = "*.example.com"
				name        = "test"
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
							data "migadu_rewrite_rule" "test" {
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
