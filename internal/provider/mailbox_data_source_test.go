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

func TestMailboxDataSource_Schema(t *testing.T) {
	ctx := context.Background()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	provider.NewMailboxDataSource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)
	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestMailboxDataSource_API_Success(t *testing.T) {
	testCases := map[string]struct {
		localPart string
		domain    string
		state     []model.Mailbox
		want      model.Mailbox
	}{
		"single": {
			localPart: "test",
			domain:    "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
				},
			},
			want: model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		"multiple": {
			localPart: "test",
			domain:    "example.com",
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
					Name:       "other",
				},
			},
			want: model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		"idna": {
			localPart: "test",
			domain:    "hoß.de",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: model.Mailbox{
				LocalPart:  "test",
				DomainName: "hoß.de",
				Address:    "test@xn--ho-hia.de",
				Name:       "Some Name",
			},
		},
		"idna-punycode": {
			localPart: "test",
			domain:    "xn--ho-hia.de",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: model.Mailbox{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Name:       "Some Name",
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
							data "migadu_mailbox" "test" {
								local_part  = "%s"
								domain_name = "%s"
							}
						`, testCase.localPart, testCase.domain),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", testCase.want.LocalPart, testCase.want.DomainName)),
							resource.TestCheckResourceAttr("data.migadu_mailbox.test", "local_part", testCase.want.LocalPart),
							resource.TestCheckResourceAttr("data.migadu_mailbox.test", "domain_name", testCase.want.DomainName),
							resource.TestCheckResourceAttr("data.migadu_mailbox.test", "address", testCase.want.Address),
							resource.TestCheckResourceAttr("data.migadu_mailbox.test", "name", testCase.want.Name),
						),
					},
				},
			})
		})
	}
}

func TestMailboxDataSource_API_Error(t *testing.T) {
	testCases := map[string]APIErrorTestCase{
		"error-401": {
			StatusCode: http.StatusUnauthorized,
			ErrorRegex: "GetMailbox: status: 401",
		},
		"error-404": {
			StatusCode: http.StatusNotFound,
			ErrorRegex: "GetMailbox: status: 404",
		},
		"error-500": {
			StatusCode: http.StatusInternalServerError,
			ErrorRegex: "GetMailbox: status: 500",
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
							data "migadu_mailbox" "test" {
								local_part  = "test"
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

func TestMailboxDataSource_Configuration_Errors(t *testing.T) {
	testCases := map[string]ConfigurationErrorTestCase{
		"empty-domain-name": {
			Configuration: `
				domain_name = ""
				local_part  = "test"
			`,
			ErrorRegex: "Attribute domain_name string length must be at least 1",
		},
		"empty-local-part": {
			Configuration: `
				domain_name = "example.com"
				local_part  = ""
			`,
			ErrorRegex: "Attribute local_part string length must be at least 1",
		},
		"missing-domain-name": {
			Configuration: `
				local_part = "test"
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
				local_part = "test"
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
							data "migadu_mailbox" "test" {
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
