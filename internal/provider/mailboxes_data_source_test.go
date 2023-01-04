//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"github.com/metio/terraform-provider-migadu/migadu/simulator"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestMailboxesDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		state  []model.Mailbox
		want   *model.Mailboxes
	}{
		{
			name:   "empty",
			domain: "example.com",
			want: &model.Mailboxes{
				Mailboxes: []model.Mailbox{},
			},
		},
		{
			name:   "single",
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "test",
				},
			},
			want: &model.Mailboxes{
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
		{
			name:   "multiple",
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
			want: &model.Mailboxes{
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
		{
			name:   "filtered",
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
			want: &model.Mailboxes{
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
		{
			name:   "idna",
			domain: "hoß.de",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: &model.Mailboxes{
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							data "migadu_mailboxes" "test" {
								domain_name = "%s"
							}
						`, tt.domain),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "mailboxes.#", fmt.Sprintf("%v", len(tt.want.Mailboxes))),
							resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "id", tt.domain),
						),
					},
				},
			})
		})
	}
}

func TestMailboxesDataSource_API_Error(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		error      string
	}{
		{
			name:       "error-401",
			domain:     "example.com",
			statusCode: http.StatusUnauthorized,
			error:      "status: 401",
		},
		{
			name:       "error-404",
			domain:     "example.com",
			statusCode: http.StatusNotFound,
			error:      "status: 404",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			statusCode: http.StatusInternalServerError,
			error:      "status: 500",
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
							data "migadu_mailboxes" "test" {
								domain_name = "%s"
							}
						`, tt.domain),
						ExpectError: regexp.MustCompile(tt.error),
					},
				},
			})
		})
	}
}

func TestMailboxesDataSource_Configuration_Errors(t *testing.T) {
	tests := []struct {
		name          string
		configuration string
		error         string
	}{
		{
			name: "empty-domain-name",
			configuration: `
				domain_name = ""
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name:          "missing-domain-name",
			configuration: ``,
			error:         `The argument "domain_name" is required, but no definition was found`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							data "migadu_mailboxes" "test" {
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
