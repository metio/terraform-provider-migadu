//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-migadu/internal/provider"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"github.com/metio/terraform-provider-migadu/migadu/simulator"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestAliasResource_API_Success(t *testing.T) {
	testCases := map[string]ResourceTestCase[model.Alias]{
		"single-destination": {
			Create: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:  "test",
					DomainName: "example.com",
					Destinations: []string{
						"other@example.com",
					},
				},
				Want: model.Alias{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"other@example.com",
					},
				},
			},
			Update: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:  "test",
					DomainName: "example.com",
					Destinations: []string{
						"another@example.com",
					},
				},
				Want: model.Alias{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"another@example.com",
					},
				},
			},
		},
		"multiple-destinations": {
			Create: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:  "test",
					DomainName: "example.com",
					Destinations: []string{
						"other@example.com",
						"some@example.com",
					},
				},
				Want: model.Alias{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"other@example.com",
						"some@example.com",
					},
				},
			},
			Update: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:  "test",
					DomainName: "example.com",
					Destinations: []string{
						"another@example.com",
						"some@example.com",
					},
				},
				Want: model.Alias{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"another@example.com",
						"some@example.com",
					},
				},
			},
		},
		"idna-domain": {
			Create: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:    "test",
					DomainName:   "hoß.de",
					Destinations: []string{"other@hoß.de"},
				},
				Want: model.Alias{
					LocalPart:  "test",
					DomainName: "hoß.de",
					Address:    "test@xn--ho-hia.de",
					Destinations: []string{
						"other@hoß.de",
					},
				},
			},
			Update: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:  "test",
					DomainName: "hoß.de",
					Destinations: []string{
						"another@hoß.de",
					},
				},
				Want: model.Alias{
					LocalPart:  "test",
					DomainName: "hoß.de",
					Address:    "test@xn--ho-hia.de",
					Destinations: []string{
						"another@hoß.de",
					},
				},
			},
			ImportIgnore: []string{"destinations"}, // ImportStateVerify does not work with SemanticEquals
		},
		"change-local-part": {
			Create: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:    "test",
					DomainName:   "example.com",
					Destinations: []string{"other@example.com"},
				},
				Want: model.Alias{
					LocalPart:    "test",
					DomainName:   "example.com",
					Address:      "test@example.com",
					Destinations: []string{"other@example.com"},
				},
			},
			Update: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:    "different",
					DomainName:   "example.com",
					Destinations: []string{"another@example.com"},
				},
				Want: model.Alias{
					LocalPart:    "different",
					DomainName:   "example.com",
					Address:      "different@example.com",
					Destinations: []string{"another@example.com"},
				},
			},
		},
		"change-domain-name": {
			Create: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:    "test",
					DomainName:   "example.com",
					Destinations: []string{"other@example.com"},
				},
				Want: model.Alias{
					LocalPart:    "test",
					DomainName:   "example.com",
					Address:      "test@example.com",
					Destinations: []string{"other@example.com"},
				},
			},
			Update: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:    "test",
					DomainName:   "different.com",
					Destinations: []string{"another@different.com"},
				},
				Want: model.Alias{
					LocalPart:    "test",
					DomainName:   "different.com",
					Address:      "test@different.com",
					Destinations: []string{"another@different.com"},
				},
			},
		},
		"change-punycode-domain-name": {
			Create: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:    "test",
					DomainName:   "xn--ho-hia.de",
					Destinations: []string{"other@xn--ho-hia.de"},
				},
				Want: model.Alias{
					LocalPart:    "test",
					DomainName:   "xn--ho-hia.de",
					Address:      "test@xn--ho-hia.de",
					Destinations: []string{"other@xn--ho-hia.de"},
				},
			},
			Update: ResourceTestStep[model.Alias]{
				Send: model.Alias{
					LocalPart:    "test",
					DomainName:   "hoß.de",
					Destinations: []string{"another@hoß.de"},
				},
				Want: model.Alias{
					LocalPart:    "test",
					DomainName:   "hoß.de",
					Address:      "test@xn--ho-hia.de",
					Destinations: []string{"another@hoß.de"},
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
							resource "migadu_alias" "test" {
								local_part   = "%s"
								domain_name  = "%s"
								destinations = %s
							}
						`, testCase.Create.Send.LocalPart, testCase.Create.Send.DomainName, strings.ReplaceAll(fmt.Sprintf("%+q", testCase.Create.Send.Destinations), "\" \"", "\",\"")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_alias.test", "id", provider.CreateAliasIDString(testCase.Create.Want.LocalPart, testCase.Create.Want.DomainName)),
							resource.TestCheckResourceAttr("migadu_alias.test", "local_part", testCase.Create.Want.LocalPart),
							resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", testCase.Create.Want.DomainName),
							resource.TestCheckResourceAttr("migadu_alias.test", "address", testCase.Create.Want.Address),
							resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", fmt.Sprintf("%v", len(testCase.Create.Want.Destinations))),
							resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", testCase.Create.Want.Destinations[0]),
							resource.TestCheckResourceAttr("migadu_alias.test", "is_internal", fmt.Sprintf("%v", testCase.Create.Want.IsInternal)),
							resource.TestCheckResourceAttr("migadu_alias.test", "expirable", fmt.Sprintf("%v", testCase.Create.Want.Expirable)),
							resource.TestCheckResourceAttr("migadu_alias.test", "expires_on", testCase.Create.Want.ExpiresOn),
							resource.TestCheckResourceAttr("migadu_alias.test", "remove_upon_expiry", fmt.Sprintf("%v", testCase.Create.Want.RemoveUponExpiry)),
						),
					},
					{
						ResourceName:            "migadu_alias.test",
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: testCase.ImportIgnore,
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_alias" "test" {
								local_part   = "%s"
								domain_name  = "%s"
								destinations = %s
							}
						`, testCase.Update.Send.LocalPart, testCase.Update.Send.DomainName, strings.ReplaceAll(fmt.Sprintf("%+q", testCase.Update.Send.Destinations), "\" \"", "\",\"")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_alias.test", "id", provider.CreateAliasIDString(testCase.Update.Want.LocalPart, testCase.Update.Want.DomainName)),
							resource.TestCheckResourceAttr("migadu_alias.test", "local_part", testCase.Update.Want.LocalPart),
							resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", testCase.Update.Want.DomainName),
							resource.TestCheckResourceAttr("migadu_alias.test", "address", testCase.Update.Want.Address),
							resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", fmt.Sprintf("%v", len(testCase.Update.Want.Destinations))),
							resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", testCase.Update.Want.Destinations[0]),
							resource.TestCheckResourceAttr("migadu_alias.test", "is_internal", fmt.Sprintf("%v", testCase.Update.Want.IsInternal)),
							resource.TestCheckResourceAttr("migadu_alias.test", "expirable", fmt.Sprintf("%v", testCase.Update.Want.Expirable)),
							resource.TestCheckResourceAttr("migadu_alias.test", "expires_on", testCase.Update.Want.ExpiresOn),
							resource.TestCheckResourceAttr("migadu_alias.test", "remove_upon_expiry", fmt.Sprintf("%v", testCase.Update.Want.RemoveUponExpiry)),
						),
					},
				},
			})
		})
	}
}

func TestAliasResource_API_Errors(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "CreateAlias: status: 404",
		},
		"error-409": {
			StatusCode: http.StatusConflict,
			ErrorRegex: "CreateAlias: status: 409",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "CreateAlias: status: 500",
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
							resource "migadu_alias" "test" {
								local_part   = "test"
								domain_name  = "example.com"
								destinations = ["other@example.com"]
							}
						`,
						ExpectError: regexp.MustCompile(testCase.ErrorRegex),
					},
				},
			})
		})
	}
}

func TestAliasResource_Configuration_Errors(t *testing.T) {
	testCases := map[string]ConfigurationErrorTestCase{
		"empty-domain-name": {
			Configuration: `
				local_part   = "test"
				domain_name  = ""
				destinations = ["someone@example.com"]
			`,
			ErrorRegex: "Attribute domain_name string length must be at least 1",
		},
		"empty-local-part": {
			Configuration: `
				local_part   = ""
				domain_name  = "example.com"
				destinations = ["someone@example.com"]
			`,
			ErrorRegex: "Attribute local_part string length must be at least 1",
		},
		"missing-domain-name": {
			Configuration: `
				local_part   = "test"
				destinations = ["someone@example.com"]
			`,
			ErrorRegex: `The argument "domain_name" is required, but no definition was found`,
		},
		"missing-local-part": {
			Configuration: `
				domain_name  = "example.com"
				destinations = ["someone@example.com"]
			`,
			ErrorRegex: `The argument "local_part" is required, but no definition was found`,
		},
		"missing-destinations": {
			Configuration: `
				local_part = "test"
				domain_name = "example.com"
			`,
			ErrorRegex: `The argument "destinations" is required, but no definition was found`,
		},
		"empty-destinations": {
			Configuration: `
				local_part   = "test"
				domain_name  = "example.com"
				destinations = []
			`,
			ErrorRegex: `Attribute destinations set must contain at least 1 elements`,
		},
		"wrong-email-format": {
			Configuration: `
				local_part   = "test"
				domain_name  = "example.com"
				destinations = ["someone"]
			`,
			ErrorRegex: `An email must match the format 'local_part@domain'`,
		},
		"duplicate-emails": {
			Configuration: `
				local_part   = "test"
				domain_name  = "example.com"
				destinations = ["someone@hoß.de", "someone@xn--ho-hia.de"]
			`,
			ErrorRegex: `This attribute contains duplicate values of: "someone@xn--ho-hia.de"`,
		},
		"invalid-domain-name": {
			Configuration: `
				local_part   = "test"
				domain_name  = "*.example.com"
				destinations = ["someone@example.com"]
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
							resource "migadu_alias" "test" {
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
