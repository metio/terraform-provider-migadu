//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"context"
	"fmt"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-migadu/internal/provider"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"github.com/metio/terraform-provider-migadu/migadu/simulator"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestRewriteRuleResource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	provider.NewRewriteRuleResource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestRewriteRuleResource_API_Success(t *testing.T) {
	testCases := map[string]ResourceTestCase[model.RewriteRule]{
		"single": {
			Create: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      2,
					Destinations: []string{
						"security@example.com",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      2,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			Update: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "security-*",
					OrderNum:      5,
					Destinations: []string{
						"security@example.com",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "security-*",
					OrderNum:      5,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
		},
		"multiple": {
			Create: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      3,
					Destinations: []string{
						"security@example.com",
						"another@example.com",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      3,
					Destinations: []string{
						"another@example.com",
						"security@example.com",
					},
				},
			},
			Update: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "security-*",
					OrderNum:      1,
					Destinations: []string{
						"security@example.com",
						"another@example.com",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "security-*",
					OrderNum:      1,
					Destinations: []string{
						"another@example.com",
						"security@example.com",
					},
				},
			},
		},
		"idna": {
			Create: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "hoß.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@hoß.de",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "hoß.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@hoß.de",
					},
				},
			},
			Update: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "hoß.de",
					Name:          "sec",
					LocalPartRule: "security-*",
					OrderNum:      7,
					Destinations: []string{
						"security@hoß.de",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "hoß.de",
					Name:          "sec",
					LocalPartRule: "security-*",
					OrderNum:      7,
					Destinations: []string{
						"security@hoß.de",
					},
				},
			},
			ImportIgnore: []string{"destinations"},
		},
		"change-name": {
			Create: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      2,
					Destinations: []string{
						"security@example.com",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      2,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			Update: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "different",
					LocalPartRule: "sec-*",
					OrderNum:      3,
					Destinations: []string{
						"security@example.com",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "different",
					LocalPartRule: "sec-*",
					OrderNum:      3,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
		},
		"change-domain-name": {
			Create: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      6,
					Destinations: []string{
						"security@example.com",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      6,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			Update: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "different.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      9,
					Destinations: []string{
						"security@different.de",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "different.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      9,
					Destinations: []string{
						"security@different.de",
					},
				},
			},
		},
		"change-punycode-domain-name": {
			Create: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "xn--ho-hia.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      8,
					Destinations: []string{
						"security@xn--ho-hia.de",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "xn--ho-hia.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      8,
					Destinations: []string{
						"security@xn--ho-hia.de",
					},
				},
			},
			Update: ResourceTestStep[model.RewriteRule]{
				Send: model.RewriteRule{
					DomainName:    "hoß.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      8,
					Destinations: []string{
						"security@hoß.de",
					},
				},
				Want: model.RewriteRule{
					DomainName:    "hoß.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      8,
					Destinations: []string{
						"security@hoß.de",
					},
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
							resource "migadu_rewrite_rule" "test" {
								domain_name     = "%s"
								name            = "%s"
								local_part_rule = "%s"
								order_num       = %d
								destinations    = %s
							}
						`, testCase.Create.Send.DomainName, testCase.Create.Send.Name, testCase.Create.Send.LocalPartRule, testCase.Create.Send.OrderNum, strings.ReplaceAll(fmt.Sprintf("%+q", testCase.Create.Send.Destinations), "\" \"", "\",\"")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "id", fmt.Sprintf("%s/%s", testCase.Create.Want.DomainName, testCase.Create.Want.Name)),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "domain_name", testCase.Create.Want.DomainName),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "name", testCase.Create.Want.Name),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "local_part_rule", testCase.Create.Want.LocalPartRule),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "order_num", fmt.Sprintf("%d", testCase.Create.Want.OrderNum)),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "destinations.#", fmt.Sprintf("%d", len(testCase.Create.Want.Destinations))),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "destinations.0", testCase.Create.Want.Destinations[0]),
						),
					},
					{
						ResourceName:            "migadu_rewrite_rule.test",
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: testCase.ImportIgnore,
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_rewrite_rule" "test" {
								domain_name     = "%s"
								name            = "%s"
								local_part_rule = "%s"
								order_num       = %d
								destinations    = %s
							}
						`, testCase.Update.Send.DomainName, testCase.Update.Send.Name, testCase.Update.Send.LocalPartRule, testCase.Update.Send.OrderNum, strings.ReplaceAll(fmt.Sprintf("%+q", testCase.Update.Send.Destinations), "\" \"", "\",\"")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "id", fmt.Sprintf("%s/%s", testCase.Update.Want.DomainName, testCase.Update.Want.Name)),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "domain_name", testCase.Update.Want.DomainName),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "name", testCase.Update.Want.Name),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "local_part_rule", testCase.Update.Want.LocalPartRule),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "order_num", fmt.Sprintf("%d", testCase.Update.Want.OrderNum)),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "destinations.#", fmt.Sprintf("%d", len(testCase.Update.Want.Destinations))),
							resource.TestCheckResourceAttr("migadu_rewrite_rule.test", "destinations.0", testCase.Update.Want.Destinations[0]),
						),
					},
				},
			})
		})
	}
}

func TestRewriteRuleResource_API_Errors(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-409": {
			StatusCode: http.StatusConflict,
			ErrorRegex: "CreateRewriteRule: status: 409",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "CreateRewriteRule: status: 500",
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
							resource "migadu_rewrite_rule" "test" {
								domain_name     = "example.com"
								name            = "sec"
								local_part_rule = "sec-*"
								destinations    = ["security@example.com"]
							}
						`,
						ExpectError: regexp.MustCompile(testCase.ErrorRegex),
					},
				},
			})
		})
	}
}

func TestRewriteRuleResource_Configuration_Errors(t *testing.T) {
	testCases := map[string]ConfigurationErrorTestCase{
		"empty-domain-name": {
			Configuration: `
				domain_name     = ""
				name            = "test"
				local_part_rule = "prefix-*"
				destinations    = ["test@example.com"]
			`,
			ErrorRegex: "Attribute domain_name string length must be at least 1",
		},
		"missing-domain-name": {
			Configuration: `
				name            = "test"
				local_part_rule = "prefix-*"
				destinations    = ["test@example.com"]
			`,
			ErrorRegex: `The argument "domain_name" is required, but no definition was found`,
		},
		"empty-name": {
			Configuration: `
				domain_name     = "example.com"
				name            = ""
				local_part_rule = "prefix-*"
				destinations    = ["test@example.com"]
			`,
			ErrorRegex: "Attribute name string length must be at least 1",
		},
		"missing-name": {
			Configuration: `
				domain_name     = "example.com"
				local_part_rule = "prefix-*"
				destinations    = ["test@example.com"]
			`,
			ErrorRegex: `The argument "name" is required, but no definition was found`,
		},
		"empty-local-part-rule": {
			Configuration: `
				domain_name     = "example.com"
				name            = "test"
				local_part_rule = ""
				destinations    = ["test@example.com"]
			`,
			ErrorRegex: "Attribute local_part_rule string length must be at least 1",
		},
		"missing-local-part-rule": {
			Configuration: `
				domain_name     = "example.com"
				name            = "test"
				destinations    = ["test@example.com"]
			`,
			ErrorRegex: `The argument "local_part_rule" is required, but no definition was found`,
		},
		"empty-destinations": {
			Configuration: `
				domain_name     = "example.com"
				name            = "test"
				local_part_rule = "prefix-*"
				destinations    = []
			`,
			ErrorRegex: `Attribute destinations set must contain at least 1 elements`,
		},
		"wrong-email-format": {
			Configuration: `
				domain_name     = "example.com"
				name            = "test"
				local_part_rule = "prefix-*"
				destinations    = ["someone"]
			`,
			ErrorRegex: `An email must match the format 'local_part@domain'`,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							resource "migadu_rewrite_rule" "test" {
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
