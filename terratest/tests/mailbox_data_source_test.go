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

func TestMailboxDataSource(t *testing.T) {
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
			localPart: "some",
			state: []model.Mailbox{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Delegations:      []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Mailbox{
				LocalPart:        "some",
				DomainName:       "example.com",
				Address:          "some@example.com",
				Delegations:      []string{"other@example"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "some",
			state: []model.Mailbox{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Delegations:      []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
				{
					LocalPart:        "other",
					DomainName:       "example.com",
					Address:          "other@example.com",
					Delegations:      []string{"different@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Mailbox{
				LocalPart:        "some",
				DomainName:       "example.com",
				Address:          "some@example.com",
				Delegations:      []string{"other@example"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:        "test",
					DomainName:       "xn--ho-hia.de",
					Address:          "test@xn--ho-hia.de",
					Delegations:      []string{"other@xn--ho-hia.de"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Mailbox{
				LocalPart:        "test",
				DomainName:       "xn--ho-hia.de",
				Address:          "test@xn--ho-hia.de",
				Delegations:      []string{"other@xn--ho-hia.de"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state}))
			defer server.Close()

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: "../data-sources/migadu_mailbox",
				Vars: map[string]interface{}{
					"endpoint":    server.URL,
					"domain_name": tt.domain,
					"local_part":  tt.localPart,
				},
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApplyAndIdempotent(t, terraformOptions)

			assert.Equal(t, fmt.Sprintf("%s@%s", tt.localPart, tt.domain), terraform.Output(t, terraformOptions, "id"), "id")
			assert.Equal(t, tt.domain, terraform.Output(t, terraformOptions, "domain_name"), "domain_name")
			assert.Equal(t, tt.localPart, terraform.Output(t, terraformOptions, "local_part"), "local_part")
			assert.Equal(t, tt.want.Address, terraform.Output(t, terraformOptions, "address"), "address")
			assert.Equal(t, tt.want.ExpiresOn, terraform.Output(t, terraformOptions, "expires_on"), "expires_on")
		})
	}
}
