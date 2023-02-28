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

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: "../data-sources/migadu_rewrite",
				Vars: map[string]interface{}{
					"endpoint":    server.URL,
					"domain_name": tt.domain,
					"name":        tt.slug,
				},
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApplyAndIdempotent(t, terraformOptions)

			assert.Equal(t, fmt.Sprintf("%s/%s", tt.domain, tt.slug), terraform.Output(t, terraformOptions, "id"), "id")
			assert.Equal(t, tt.domain, terraform.Output(t, terraformOptions, "domain_name"), "domain_name")
			assert.Equal(t, tt.slug, terraform.Output(t, terraformOptions, "name"), "name")
			assert.Equal(t, tt.want.LocalPartRule, terraform.Output(t, terraformOptions, "local_part_rule"), "local_part_rule")
			assert.Equal(t, fmt.Sprintf("%v", tt.want.OrderNum), terraform.Output(t, terraformOptions, "order_num"), "order_num")
		})
	}
}
