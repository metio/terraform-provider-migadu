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
	"strings"
	"testing"
)

func TestRewriteResource_API_Success(t *testing.T) {
	tests := []struct {
		name          string
		domain        string
		slug          string
		localPartRule string
		destination   string
		state         []model.Rewrite
		send          *model.Rewrite
		updatedRule   string
		want          *model.Rewrite
	}{
		{
			name:   "single",
			domain: "example.com",
			state:  []model.Rewrite{},
			send: &model.Rewrite{
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
			updatedRule: "security-*",
		},
		{
			name:   "multiple",
			domain: "example.com",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "existing",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			send: &model.Rewrite{
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
					"another@example.com",
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
					"another@example.com",
				},
			},
			updatedRule: "security-*",
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
							resource "migadu_rewrite" "test" {
								domain_name     = "%s"
								name            = "%s"
								local_part_rule = "%s"
								destinations    = %s
							}
						`, tt.domain, tt.send.Name, tt.send.LocalPartRule, strings.ReplaceAll(fmt.Sprintf("%+q", tt.send.Destinations), "\" \"", "\",\"")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_rewrite.test", "domain_name", tt.want.DomainName),
							resource.TestCheckResourceAttr("migadu_rewrite.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_rewrite.test", "local_part_rule", tt.want.LocalPartRule),
							resource.TestCheckResourceAttr("migadu_rewrite.test", "destinations.#", fmt.Sprintf("%d", len(tt.want.Destinations))),
							resource.TestCheckResourceAttr("migadu_rewrite.test", "id", fmt.Sprintf("%s/%s", tt.want.DomainName, tt.want.Name)),
						),
					},
					{
						ResourceName:      "migadu_rewrite.test",
						ImportState:       true,
						ImportStateVerify: true,
					},
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_rewrite" "test" {
								domain_name     = "%s"
								name            = "%s"
								local_part_rule = "%s"
								destinations    = %s
							}
						`, tt.domain, tt.send.Name, tt.updatedRule, strings.ReplaceAll(fmt.Sprintf("%+q", tt.send.Destinations), "\" \"", "\",\"")),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("migadu_rewrite.test", "domain_name", tt.want.DomainName),
							resource.TestCheckResourceAttr("migadu_rewrite.test", "name", tt.want.Name),
							resource.TestCheckResourceAttr("migadu_rewrite.test", "local_part_rule", tt.updatedRule),
							resource.TestCheckResourceAttr("migadu_rewrite.test", "destinations.#", fmt.Sprintf("%d", len(tt.want.Destinations))),
							resource.TestCheckResourceAttr("migadu_rewrite.test", "id", fmt.Sprintf("%s/%s", tt.want.DomainName, tt.want.Name)),
						),
					},
				},
			})
		})
	}
}

func TestRewriteResource_API_Errors(t *testing.T) {
	tests := []struct {
		name          string
		domain        string
		slug          string
		localPartRule string
		destination   string
		statusCode    int
		state         []model.Rewrite
		error         string
	}{
		{
			name:          "error-400",
			domain:        "example.com",
			slug:          "sec",
			localPartRule: "sec-*",
			destination:   "security@example.com",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			error: "CreateRewrite: status: 400",
		},
		{
			name:          "error-500",
			domain:        "example.com",
			slug:          "sec",
			localPartRule: "sec-*",
			destination:   "security@example.com",
			statusCode:    http.StatusInternalServerError,
			error:         "CreateRewrite: status: 500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig(server.URL) + fmt.Sprintf(`
							resource "migadu_rewrite" "test" {
								domain_name     = "%s"
								name            = "%s"
								local_part_rule = "%s"
								destinations    = ["%s"]
							}
						`, tt.domain, tt.slug, tt.localPartRule, tt.destination),
						ExpectError: regexp.MustCompile(tt.error),
					},
				},
			})
		})
	}
}

func TestRewriteResource_Configuration_Errors(t *testing.T) {
	tests := []struct {
		name          string
		configuration string
		error         string
	}{
		{
			name: "empty-domain-name",
			configuration: `
				domain_name     = ""
				name            = "test"
				local_part_rule = "prefix-*"
				destinations    = ["test@example.com"]
			`,
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name: "missing-domain-name",
			configuration: `
				name            = "test"
				local_part_rule = "prefix-*"
				destinations    = ["test@example.com"]
			`,
			error: `The argument "domain_name" is required, but no definition was found`,
		},
		{
			name: "empty-name",
			configuration: `
				domain_name     = "example.com"
				name            = ""
				local_part_rule = "prefix-*"
				destinations    = ["test@example.com"]
			`,
			error: "Attribute name string length must be at least 1",
		},
		{
			name: "missing-name",
			configuration: `
				domain_name     = "example.com"
				local_part_rule = "prefix-*"
				destinations    = ["test@example.com"]
			`,
			error: `The argument "name" is required, but no definition was found`,
		},
		{
			name: "empty-local-part-rule",
			configuration: `
				domain_name     = "example.com"
				name            = "test"
				local_part_rule = ""
				destinations    = ["test@example.com"]
			`,
			error: "Attribute local_part_rule string length must be at least 1",
		},
		{
			name: "missing-local-part-rule",
			configuration: `
				domain_name     = "example.com"
				name            = "test"
				destinations    = ["test@example.com"]
			`,
			error: `The argument "local_part_rule" is required, but no definition was found`,
		},
		{
			name: "empty-destinations",
			configuration: `
				domain_name     = "example.com"
				name            = "test"
				local_part_rule = "prefix-*"
				destinations    = []
			`,
			error: `Attribute destinations list must contain at least 1 elements`,
		},
		{
			name: "missing-destinations",
			configuration: `
				domain_name     = "example.com"
				name            = "test"
				local_part_rule = "prefix-*"
			`,
			error: `No attribute specified when one \(and only one\) of \[destinations\] is required`,
		},
		{
			name: "empty-destinations-punycode",
			configuration: `
				domain_name           = "example.com"
				name                  = "test"
				local_part_rule       = "prefix-*"
				destinations_punycode = []
			`,
			error: `Attribute destinations_punycode list must contain at least 1 elements`,
		},
		{
			name: "missing-destinations-punycode",
			configuration: `
				domain_name           = "example.com"
				name                  = "test"
				local_part_rule       = "prefix-*"
			`,
			error: `(?s)No attribute specified when one \(and only one\) of \[destinations_punycode\] is(.*)required`,
		},
		{
			name: "multiple-destination-attributes",
			configuration: `
				domain_name           = "example.com"
				name                  = "test"
				local_part_rule       = "prefix-*"
				destinations          = ["test@example.com"]
				destinations_punycode = ["test@example.com"]
			`,
			error: `2 attributes specified when one \(and only one\) of \[destinations\] is required`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: providerConfig("https://localhost:12345") + fmt.Sprintf(`
							resource "migadu_rewrite" "test" {
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
