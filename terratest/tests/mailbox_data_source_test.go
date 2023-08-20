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
	testCases := map[string]struct {
		localPart string
		domain    string
		state     []model.Mailbox
		want      model.Mailbox
	}{
		"single": {
			localPart: "some",
			domain:    "example.com",
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
			want: model.Mailbox{
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
		"multiple": {
			localPart: "some",
			domain:    "example.com",
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
			want: model.Mailbox{
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
		"idna": {
			localPart: "test",
			domain:    "hoß.de",
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
			want: model.Mailbox{
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
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: testCase.state}))
			defer server.Close()

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: "../data-sources/migadu_mailbox",
				Vars: map[string]interface{}{
					"endpoint":    server.URL,
					"local_part":  testCase.localPart,
					"domain_name": testCase.domain,
				},
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApplyAndIdempotent(t, terraformOptions)

			assert.Equal(t, fmt.Sprintf("%s@%s", testCase.localPart, testCase.domain), terraform.Output(t, terraformOptions, "id"), "id")
			assert.Equal(t, testCase.localPart, terraform.Output(t, terraformOptions, "local_part"), "local_part")
			assert.Equal(t, testCase.domain, terraform.Output(t, terraformOptions, "domain_name"), "domain_name")
			assert.Equal(t, testCase.want.Address, terraform.Output(t, terraformOptions, "address"), "address")
			assert.Equal(t, testCase.want.ExpiresOn, terraform.Output(t, terraformOptions, "expires_on"), "expires_on")
		})
	}
}
