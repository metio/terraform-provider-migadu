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

func TestRewritesDataSource_API_Success(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		state  []model.Rewrite
		want   *model.Rewrites
	}{
		{
			name:   "empty",
			domain: "example.com",
			want: &model.Rewrites{
				Rewrites: []model.Rewrite{},
			},
		},
		{
			name:   "single",
			domain: "example.com",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@example.com",
					},
				},
			},
			want: &model.Rewrites{
				Rewrites: []model.Rewrite{
					{
						DomainName:    "example.com",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest@example.com",
						},
					},
				},
			},
		},
		{
			name:   "multiple",
			domain: "example.com",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest1@example.com",
						"dest2@example.com",
					},
				},
				{
					DomainName:    "example.com",
					Name:          "Other Rule",
					LocalPartRule: "*-suffix",
					OrderNum:      1,
					Destinations: []string{
						"dest3@example.com",
						"dest4@example.com",
					},
				},
			},
			want: &model.Rewrites{
				Rewrites: []model.Rewrite{
					{
						DomainName:    "example.com",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest1@example.com",
							"dest2@example.com",
						},
					},
					{
						DomainName:    "example.com",
						Name:          "Other Rule",
						LocalPartRule: "*-suffix",
						OrderNum:      1,
						Destinations: []string{
							"dest3@example.com",
							"dest4@example.com",
						},
					},
				},
			},
		},
		{
			name:   "filtered",
			domain: "example.com",
			state: []model.Rewrite{
				{
					DomainName:    "different.com",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest1@different.com",
						"dest2@different.com",
					},
				},
				{
					DomainName:    "example.com",
					Name:          "Other Rule",
					LocalPartRule: "*-suffix",
					OrderNum:      1,
					Destinations: []string{
						"dest3@example.com",
						"dest4@example.com",
					},
				},
			},
			want: &model.Rewrites{
				Rewrites: []model.Rewrite{
					{
						DomainName:    "example.com",
						Name:          "Other Rule",
						LocalPartRule: "*-suffix",
						OrderNum:      1,
						Destinations: []string{
							"dest3@example.com",
							"dest4@example.com",
						},
					},
				},
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			state: []model.Rewrite{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "Some Rule",
					LocalPartRule: "prefix-*",
					OrderNum:      0,
					Destinations: []string{
						"dest@xn--ho-hia.de",
					},
				},
			},
			want: &model.Rewrites{
				Rewrites: []model.Rewrite{
					{
						DomainName:    "xn--ho-hia.de",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest@xn--ho-hia.de",
						},
					},
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
							data "migadu_rewrites" "test" {
								domain_name = "%s"
							}
						`, tt.domain),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.migadu_rewrites.test", "domain_name", tt.domain),
							resource.TestCheckResourceAttr("data.migadu_rewrites.test", "rewrites.#", fmt.Sprintf("%v", len(tt.want.Rewrites))),
							resource.TestCheckResourceAttr("data.migadu_rewrites.test", "id", tt.domain),
						),
					},
				},
			})
		})
	}
}

func TestRewritesDataSource_API_Errors(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		error      string
	}{
		{
			name:       "error-404",
			domain:     "example.com",
			statusCode: http.StatusNotFound,
			error:      "GetRewrites: status: 404",
		},
		{
			name:       "error-500",
			domain:     "example.com",
			statusCode: http.StatusInternalServerError,
			error:      "GetRewrites: status: 500",
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
							data "migadu_rewrites" "test" {
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

func TestRewritesDataSource_Configuration_Errors(t *testing.T) {
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
							data "migadu_rewrites" "test" {
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
