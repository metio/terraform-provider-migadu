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

func TestRewriteDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		slug   string
		state  []model.Rewrite
		want   *model.Rewrite
	}{
		{
			name:   "single",
			domain: "example.com",
			slug:   "test",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@example.com",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "test",
				LocalPartRule: "prefix-*",
				OrderNum:      0,
				Destinations: []string{
					"dest@example.com",
				},
			},
		},
		{
			name:   "multiple",
			domain: "example.com",
			slug:   "test",
			state: []model.Rewrite{
				{
					DomainName:    "different.com",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@different.com",
					},
				},
				{
					DomainName:    "example.com",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@example.com",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "test",
				LocalPartRule: "prefix-*",
				OrderNum:      0,
				Destinations: []string{
					"dest@example.com",
				},
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			slug:   "test",
			state: []model.Rewrite{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "test",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@xn--ho-hia.de",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "xn--ho-hia.de",
				Name:          "test",
				LocalPartRule: "prefix-*",
				OrderNum:      0,
				Destinations: []string{
					"dest@xn--ho-hia.de",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: tt.state}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							data "migadu_rewrite" "test" {
								domain_name = "%s"
								name        = "%s"
							}
						`, tt.domain, tt.slug),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_rewrite.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("data.migadu_rewrite.test", "name", tt.slug),
							resource.TestCheckResourceAttr("data.migadu_rewrite.test", "id", fmt.Sprintf("%s/%s", tt.domain, tt.slug)),
						),
					},
				},
			})
		})
	}
}

func TestRewriteDataSource_API_Errors(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		slug       string
		statusCode int
		error      string
	}{
		{
			name:       "error-404",
			domain:     "example.com",
			slug:       "some",
			statusCode: http.StatusNotFound,
			error:      "GetRewrite: status: 404",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			slug:       "some",
			statusCode: http.StatusInternalServerError,
			error:      "GetRewrite: status: 500",
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
							data "migadu_rewrite" "test" {
								domain_name = "%s"
								name        = "%s"
							}
						`, tt.domain, tt.slug),
						ExpectError: regexp.MustCompile(tt.error),
					},
				},
			})
		})
	}
}

func TestRewriteDataSource_Configuration_Errors(t *testing.T) {
	tests := []struct {
		name          string
		configuration string
		error         string
	}{
		{
			name: "empty-domain-name",
			configuration: `
				domain_name = ""
				name        = "slug"
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name: "empty-name",
			configuration: `
				domain_name = "example.com"
				name        = ""
			`,
			error: "Attribute name string length must be at least 1",
		},
		{
			name: "missing-domain-name",
			configuration: `
				name = "slug"
			`,
			error: `The argument "domain_name" is required, but no definition was found`,
		},
		{
			name: "missing-name",
			configuration: `
				domain_name = "example.com"
			`,
			error: `The argument "name" is required, but no definition was found`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							data "migadu_rewrite" "test" {
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
