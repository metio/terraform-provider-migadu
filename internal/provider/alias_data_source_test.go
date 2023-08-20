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

func TestAliasDataSource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	provider.NewAliasDataSource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestAliasDataSource_API_Success(t *testing.T) {
	testCases := map[string]struct {
		localPart string
		domain    string
		state     []model.Alias
		want      model.Alias
	}{
		"single": {
			localPart: "some",
			domain:    "example.com",
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
			want: model.Alias{
				Address:          "some@example.com",
				Destinations:     []string{"other@example"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		"multiple": {
			localPart: "some",
			domain:    "example.com",
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
			want: model.Alias{
				Address:          "some@example.com",
				Destinations:     []string{"other@example"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		"idna": {
			localPart: "test",
			domain:    "hoß.de",
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
			want: model.Alias{
				Address:          "test@xn--ho-hia.de",
				Destinations:     []string{"other@xn--ho-hia.de"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		"idna-punycode": {
			localPart: "test",
			domain:    "xn--ho-hia.de",
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
			want: model.Alias{
				Address:          "test@xn--ho-hia.de",
				Destinations:     []string{"other@xn--ho-hia.de"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
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
							data "migadu_alias" "test" {
								local_part  = "%s"
								domain_name = "%s"
							}
						`, testCase.localPart, testCase.domain),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_alias.test", "id", provider.CreateAliasIDString(testCase.localPart, testCase.domain)),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "local_part", testCase.localPart),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "domain_name", testCase.domain),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "address", testCase.want.Address),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "destinations.#", fmt.Sprintf("%v", len(testCase.want.Destinations))),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "destinations.0", testCase.want.Destinations[0]),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "is_internal", fmt.Sprintf("%v", testCase.want.IsInternal)),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "expirable", fmt.Sprintf("%v", testCase.want.Expirable)),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "expires_on", testCase.want.ExpiresOn),
							resource.TestCheckResourceAttr("data.migadu_alias.test", "remove_upon_expiry", fmt.Sprintf("%v", testCase.want.RemoveUponExpiry)),
						),
					},
				},
			})
		})
	}
}

func TestAliasDataSource_API_Errors(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "GetAlias: status: 404",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "GetAlias: status: 500",
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
							data "migadu_alias" "test" {
								local_part  = "other"
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

func TestAliasDataSource_Configuration_Errors(t *testing.T) {
	testCases := map[string]ConfigurationErrorTestCase{
		"empty-domain-name": {
			Configuration: `
				local_part  = "test"
				domain_name = ""
			`,
			ErrorRegex: "Attribute domain_name string length must be at least 1",
		},
		"empty-local-part": {
			Configuration: `
				local_part  = ""
				domain_name = "example.com"
			`,
			ErrorRegex: "Attribute local_part string length must be at least 1",
		},
		"missing-domain-name": {
			Configuration: `
				local_part  = "test"
			`,
			ErrorRegex: `The argument "domain_name" is required, but no definition was found`,
		},
		"missing-local-part": {
			Configuration: `
				domain_name = "example.com"
			`,
			ErrorRegex: `The argument "local_part" is required, but no definition was found`,
		},
		"invalid-domain-name": {
			Configuration: `
				local_part  = "test"
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
							data "migadu_alias" "test" {
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
