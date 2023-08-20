//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package acceptance_test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"github.com/metio/terraform-provider-migadu/migadu/simulator"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestRewriteDataSource(t *testing.T) {
	testCases := map[string]struct {
		domain string
		name   string
		state  []model.RewriteRule
		want   model.RewriteRule
	}{
		"single": {
			domain: "example.com",
			name:   "test",
			state: []model.RewriteRule{
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
			want: model.RewriteRule{
				DomainName:    "example.com",
				Name:          "test",
				LocalPartRule: "prefix-*",
				OrderNum:      0,
				Destinations: []string{
					"dest@example.com",
				},
			},
		},
		"multiple": {
			domain: "example.com",
			name:   "test",
			state: []model.RewriteRule{
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
			want: model.RewriteRule{
				DomainName:    "example.com",
				Name:          "test",
				LocalPartRule: "prefix-*",
				OrderNum:      0,
				Destinations: []string{
					"dest@example.com",
				},
			},
		},
		"idna": {
			domain: "hoß.de",
			name:   "test",
			state: []model.RewriteRule{
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
			want: model.RewriteRule{
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
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: testCase.state}))
			defer server.Close()

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: "../data-sources/migadu_rewrite_rule",
				Vars: map[string]interface{}{
					"endpoint":    server.URL,
					"domain_name": testCase.domain,
					"name":        testCase.name,
				},
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApplyAndIdempotent(t, terraformOptions)

			assert.Equal(t, fmt.Sprintf("%s/%s", testCase.domain, testCase.name), terraform.Output(t, terraformOptions, "id"), "id")
			assert.Equal(t, testCase.domain, terraform.Output(t, terraformOptions, "domain_name"), "domain_name")
			assert.Equal(t, testCase.name, terraform.Output(t, terraformOptions, "name"), "name")
			assert.Equal(t, testCase.want.LocalPartRule, terraform.Output(t, terraformOptions, "local_part_rule"), "local_part_rule")
			assert.Equal(t, fmt.Sprintf("%v", testCase.want.OrderNum), terraform.Output(t, terraformOptions, "order_num"), "order_num")
		})
	}
}
