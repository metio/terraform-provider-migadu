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

func TestIdentitiesDataSource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	provider.NewIdentitiesDataSource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestIdentitiesDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name      string
		domain    string
		localPart string
		state     []model.Identity
		want      model.Identities
	}{
		{
			name:      "empty",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Identity{},
			want:      model.Identities{},
		},
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:            "other",
					DomainName:           "example.com",
					Address:              "other@example.com",
					Name:                 "Some Name",
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
			want: model.Identities{
				Identities: []model.Identity{
					{
						LocalPart:            "other",
						DomainName:           "example.com",
						Address:              "other@example.com",
						Name:                 "Some Name",
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
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "another",
					DomainName: "example.com",
					Address:    "another@example.com",
					Name:       "Another Name",
				},
			},
			want: model.Identities{
				Identities: []model.Identity{
					{
						LocalPart:  "other",
						DomainName: "example.com",
						Address:    "other@example.com",
						Name:       "Some Name",
					},
					{
						LocalPart:  "another",
						DomainName: "example.com",
						Address:    "another@example.com",
						Name:       "Another Name",
					},
				},
			},
		},
		{
			name:      "filtered",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "another",
					DomainName: "different.com",
					Address:    "another@different.com",
					Name:       "Another Name",
				},
			},
			want: model.Identities{
				Identities: []model.Identity{
					{
						LocalPart:  "other",
						DomainName: "example.com",
						Address:    "other@example.com",
						Name:       "Some Name",
					},
				},
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "xn--ho-hia.de",
					Address:    "other@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: model.Identities{
				Identities: []model.Identity{
					{
						LocalPart:  "other",
						DomainName: "xn--ho-hia.de",
						Address:    "other@xn--ho-hia.de",
						Name:       "Some Name",
					},
				},
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
							data "migadu_identities" "test" {
								domain_name = "%s"
								local_part  = "%s"
							}
						`, tt.domain, tt.localPart),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_identities.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("data.migadu_identities.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("data.migadu_identities.test", "identities.#", fmt.Sprintf("%v", len(tt.want.Identities))),
							resource.TestCheckResourceAttr("data.migadu_identities.test", "id", fmt.Sprintf("%s@%s", tt.localPart, tt.domain)),
						),
					},
				},
			})
		})
	}
}

func TestIdentitiesDataSource_API_Errors(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		error      string
	}{
		{
			name:       "error-404",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusNotFound,
			error:      "GetIdentities: status: 404",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusInternalServerError,
			error:      "GetIdentities: status: 500",
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
							data "migadu_identities" "test" {
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

func TestIdentitiesDataSource_Configuration_Errors(t *testing.T) {
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
							data "migadu_identities" "test" {
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
