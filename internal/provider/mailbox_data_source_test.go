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

func TestMailboxDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name      string
		domain    string
		localPart string
		state     []model.Mailbox
		want      *model.Mailbox
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "test",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "test",
				},
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "other",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Name:       "Some Name",
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
							data "migadu_mailbox" "test" {
								domain_name = "%s"
								local_part  = "%s"
							}
						`, tt.domain, tt.localPart),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_mailbox.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("data.migadu_mailbox.test", "local_part", tt.localPart),
							resource.TestCheckResourceAttr("data.migadu_mailbox.test", "id", fmt.Sprintf("%s@%s", tt.localPart, tt.domain)),
						),
					},
				},
			})
		})
	}
}

func TestMailboxDataSource_API_Error(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		error      string
	}{
		{
			name:       "error-401",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusUnauthorized,
			error:      "status: 401",
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			error:     "status: 404",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
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
							data "migadu_mailbox" "test" {
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

func TestMailboxDataSource_Configuration_Errors(t *testing.T) {
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
							data "migadu_mailbox" "test" {
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
