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

func TestMailboxesDataSource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	provider.NewMailboxesDataSource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestMailboxesDataSource_API_Success(t *testing.T) {
	testCases := map[string]struct {
		domain string
		state  []model.Mailbox
		want   model.Mailboxes
	}{
		"empty": {
			domain: "example.com",
			want: model.Mailboxes{
				Mailboxes: []model.Mailbox{},
			},
		},
		"single": {
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "test",
				},
			},
			want: model.Mailboxes{
				Mailboxes: []model.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "example.com",
						Address:    "test@example.com",
						Name:       "Some Name",
					},
				},
			},
		},
		"multiple": {
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Other Name",
				},
			},
			want: model.Mailboxes{
				Mailboxes: []model.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "example.com",
						Address:    "test@example.com",
						Name:       "Some Name",
					},
					{
						LocalPart:  "other",
						DomainName: "example.com",
						Address:    "other@example.com",
						Name:       "Other Name",
					},
				},
			},
		},
		"filtered": {
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "different.com",
					Address:    "test@different.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Other Name",
				},
			},
			want: model.Mailboxes{
				Mailboxes: []model.Mailbox{
					{
						LocalPart:  "other",
						DomainName: "example.com",
						Address:    "other@example.com",
						Name:       "Other Name",
					},
				},
			},
		},
		"idna": {
			domain: "hoß.de",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: model.Mailboxes{
				Mailboxes: []model.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "xn--ho-hia.de",
						Address:    "test@xn--ho-hia.de",
						Name:       "Some Name",
					},
				},
			},
		},
		"idna-punycode": {
			domain: "xn--ho-hia.de",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: model.Mailboxes{
				Mailboxes: []model.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "xn--ho-hia.de",
						Address:    "test@xn--ho-hia.de",
						Name:       "Some Name",
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: testCase.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							data "migadu_mailboxes" "test" {
								domain_name = "%s"
							}
						`, testCase.domain),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "id", testCase.domain),
							resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "domain_name", testCase.domain),
							resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "mailboxes.#", fmt.Sprintf("%v", len(testCase.want.Mailboxes))),
						),
					},
				},
			})
		})
	}
}

func TestMailboxesDataSource_API_Error(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-401": {
			StatusCode: http.StatusUnauthorized,
			ErrorRegex: "GetMailboxes: status: 401",
		},
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "GetMailboxes: status: 404",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "GetMailboxes: status: 500",
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
							data "migadu_mailboxes" "test" {
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

func TestMailboxesDataSource_Configuration_Errors(t *testing.T) {
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
							data "migadu_mailboxes" "test" {
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
