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

func TestIdentityDataSource(t *testing.T) {
	tests := []struct {
		name      string
		domain    string
		localPart string
		identity  string
		state     []model.Identity
		want      *model.Identity
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			identity:  "someone",
			state: []model.Identity{
				{
					LocalPart:            "someone",
					DomainName:           "example.com",
					Address:              "someone@example.com",
					Name:                 "Some Identity",
					MaySend:              true,
					MayReceive:           true,
					MayAccessImap:        true,
					MayAccessPop3:        true,
					MayAccessManageSieve: true,
					Password:             "secret",
					FooterActive:         false,
					FooterPlainBody:      "",
					FooterHtmlBody:       "",
				},
			},
			want: &model.Identity{
				LocalPart:            "someone",
				DomainName:           "example.com",
				Address:              "someone@example.com",
				Name:                 "Some Identity",
				MaySend:              true,
				MayReceive:           true,
				MayAccessImap:        true,
				MayAccessPop3:        true,
				MayAccessManageSieve: true,
				Password:             "secret",
				FooterActive:         false,
				FooterPlainBody:      "",
				FooterHtmlBody:       "",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			identity:  "someone",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "example.com",
					Address:    "someone@example.com",
					Name:       "Some Identity",
				},
				{
					LocalPart:  "another",
					DomainName: "example.com",
					Address:    "another@example.com",
					Name:       "Another Identity",
				},
			},
			want: &model.Identity{
				LocalPart:  "someone",
				DomainName: "example.com",
				Address:    "someone@example.com",
				Name:       "Some Identity",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			identity:  "someone",
			state: []model.Identity{
				{
					LocalPart:  "someone",
					DomainName: "xn--ho-hia.de",
					Address:    "someone@xn--ho-hia.de",
					Name:       "Some Identity",
				},
			},
			want: &model.Identity{
				LocalPart:  "someone",
				DomainName: "xn--ho-hia.de",
				Address:    "someone@xn--ho-hia.de",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state}))
			defer server.Close()

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: "../data-sources/migadu_identity",
				Vars: map[string]interface{}{
					"endpoint":    server.URL,
					"domain_name": tt.domain,
					"local_part":  tt.localPart,
					"identity":    tt.identity,
				},
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApplyAndIdempotent(t, terraformOptions)

			assert.Equal(t, fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.identity), terraform.Output(t, terraformOptions, "id"), "id")
			assert.Equal(t, tt.domain, terraform.Output(t, terraformOptions, "domain_name"), "domain_name")
			assert.Equal(t, tt.localPart, terraform.Output(t, terraformOptions, "local_part"), "local_part")
			assert.Equal(t, tt.identity, terraform.Output(t, terraformOptions, "identity"), "identity")
			assert.Equal(t, tt.want.Address, terraform.Output(t, terraformOptions, "address"), "address")
		})
	}
}
