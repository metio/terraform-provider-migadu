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
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestAliasResource(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		localPart    string
		destinations []string
		want         *model.Alias
		error        string
	}{
		{
			name:      "no-domain",
			domain:    "",
			localPart: "test",
			destinations: []string{
				"other@example.com",
			},
			want:  &model.Alias{},
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name:      "no-local-part",
			domain:    "example.com",
			localPart: "",
			destinations: []string{
				"other@example.com",
			},
			want:  &model.Alias{},
			error: "Attribute local_part string length must be at least 1",
		},
		{
			name:         "no-destinations",
			domain:       "example.com",
			localPart:    "test",
			destinations: []string{},
			want:         &model.Alias{},
			error:        "Attribute destinations list must contain at least 1 elements",
		},
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			destinations: []string{
				"other@example.com",
			},
			want: &model.Alias{
				Address: "test@example.com",
				Destinations: []string{
					"other@example.com",
				},
				IsInternal:       false,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			destinations: []string{
				"other@example.com",
				"some@example.com",
			},
			want: &model.Alias{
				Address: "test@example.com",
				Destinations: []string{
					"other@example.com",
					"some@example.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{}))
			defer server.Close()

			config := providerConfig(server.URL) + fmt.Sprintf(`
					resource "migadu_alias" "test" {
						domain_name  = "%s"
						local_part   = "%s"
						destinations = %s
					}
				`, tt.domain, tt.localPart, strings.ReplaceAll(fmt.Sprintf("%+q", tt.destinations), "\" \"", "\",\""))

			if tt.error != "" {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config:      config,
							ExpectError: regexp.MustCompile(tt.error),
						},
					},
				})
			} else {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: config,
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", tt.domain),
								resource.TestCheckResourceAttr("migadu_alias.test", "local_part", tt.localPart),
								resource.TestCheckResourceAttr("migadu_alias.test", "address", tt.want.Address),
								resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", fmt.Sprintf("%v", len(tt.want.Destinations))),
								resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", tt.want.Destinations[0]),
								resource.TestCheckResourceAttr("migadu_alias.test", "id", fmt.Sprintf("%s@%s", tt.localPart, tt.domain)),
							),
						},
						{
							ResourceName:      "migadu_alias.test",
							ImportState:       true,
							ImportStateVerify: true,
						},
					},
				})
			}
		})
	}
}

func TestAliasResource_IDN_Punycode(t *testing.T) {
	server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{}))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(server.URL) + `
					resource "migadu_alias" "test" {
						domain_name           = "hoß.de"
						local_part            = "test"
						destinations_punycode = ["first@xn--ho-hia.de", "second@xn--ho-hia.de"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", "hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "local_part", "test"),
					resource.TestCheckResourceAttr("migadu_alias.test", "address", "test@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", "2"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", "first@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.1", "second@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.#", "2"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.0", "first@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.1", "second@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "id", "test@hoß.de"),
				),
			},
			{
				ResourceName:      "migadu_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig(server.URL) + `
					resource "migadu_alias" "test" {
						domain_name           = "hoß.de"
						local_part            = "test"
						destinations_punycode = ["third@xn--ho-hia.de"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", "hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "local_part", "test"),
					resource.TestCheckResourceAttr("migadu_alias.test", "address", "test@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", "1"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", "third@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.#", "1"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.0", "third@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "id", "test@hoß.de"),
				),
			},
		},
	})
}

func TestAliasResource_IDN_Unicode(t *testing.T) {
	server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{}))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(server.URL) + `
					resource "migadu_alias" "test" {
						domain_name  = "hoß.de"
						local_part   = "test"
						destinations = ["first@hoß.de", "second@hoß.de"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", "hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "local_part", "test"),
					resource.TestCheckResourceAttr("migadu_alias.test", "address", "test@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", "2"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", "first@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.1", "second@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.#", "2"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.0", "first@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.1", "second@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "id", "test@hoß.de"),
				),
			},
			{
				ResourceName:      "migadu_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig(server.URL) + `
					resource "migadu_alias" "test" {
						domain_name  = "hoß.de"
						local_part   = "test"
						destinations = ["third@hoß.de"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", "hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "local_part", "test"),
					resource.TestCheckResourceAttr("migadu_alias.test", "address", "test@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", "1"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", "third@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.#", "1"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_punycode.0", "third@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "id", "test@hoß.de"),
				),
			},
		},
	})
}
