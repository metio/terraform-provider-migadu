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

func TestIdentityDataSource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	provider.NewIdentityDataSource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestIdentityDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name      string
		domain    string
		localPart string
		identity  string
		state     []model.Identity
		want      model.Identity
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
			want: model.Identity{
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
			want: model.Identity{
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
			want: model.Identity{
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
	testCases := map[string]APIErrorTestCase{
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "GetIdentity: status: 404",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "GetIdentity: status: 500",
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
							data "migadu_identity" "test" {
								local_part  = "test"
								domain_name = "example.com"
								identity    = "identity"
							}
						`,
						ExpectError: regexp.MustCompile(testCase.ErrorRegex),
					},
				},
			})
		})
	}
}

func TestIdentityDataSource_Configuration_Errors(t *testing.T) {
	testCases := map[string]ConfigurationErrorTestCase{
		"empty-domain-name": {
			Configuration: `
				domain_name = ""
				local_part  = "test"
				identity    = "test"
			`,
			ErrorRegex: "Attribute domain_name string length must be at least 1",
		},
		"empty-local-part": {
			Configuration: `
				domain_name = "example.com"
				local_part  = ""
				identity    = "test"
			`,
			ErrorRegex: "Attribute local_part string length must be at least 1",
		},
		"empty-identity": {
			Configuration: `
				domain_name = "example.com"
				local_part  = "test"
				identity    = ""
			`,
			ErrorRegex: "Attribute identity string length must be at least 1",
		},
		"missing-domain-name": {
			Configuration: `
				local_part  = "test"
				identity    = "test"
			`,
			ErrorRegex: `The argument "domain_name" is required, but no definition was found`,
		},
		"missing-local-part": {
			Configuration: `
				domain_name = "example.com"
				identity    = "test"
			`,
			ErrorRegex: `The argument "local_part" is required, but no definition was found`,
		},
		"missing-identity": {
			Configuration: `
				domain_name = "example.com"
				local_part  = "test"
			`,
			ErrorRegex: `The argument "identity" is required, but no definition was found`,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							data "migadu_identity" "test" {
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
