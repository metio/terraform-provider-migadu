/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package acceptance_test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/metio/migadu-client.go/model"
	"github.com/metio/migadu-client.go/simulator"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestMailboxResource_Using_Password(t *testing.T) {
	testCases := map[string]struct {
		localPart string
		domain    string
		state     []model.Mailbox
		want      model.Mailbox
	}{
		"single": {
			localPart: "some",
			domain:    "example.com",
			state:     []model.Mailbox{},
			want: model.Mailbox{
				LocalPart:        "some",
				DomainName:       "example.com",
				Address:          "some@example.com",
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
					LocalPart:        "other",
					DomainName:       "example.com",
					Address:          "other@example.com",
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
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		"idna": {
			localPart: "test",
			domain:    "hoß.de",
			state:     []model.Mailbox{},
			want: model.Mailbox{
				LocalPart:        "test",
				DomainName:       "xn--ho-hia.de",
				Address:          "test@xn--ho-hia.de",
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
				TerraformDir: "../resources/migadu_mailbox/password",
				Vars: map[string]interface{}{
					"endpoint":    server.URL,
					"domain_name": testCase.domain,
					"local_part":  testCase.localPart,
				},
			})

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApplyAndIdempotent(t, terraformOptions)

			assert.Equal(t, fmt.Sprintf("%s@%s", testCase.localPart, testCase.domain), terraform.Output(t, terraformOptions, "id"), "id")
			assert.Equal(t, testCase.domain, terraform.Output(t, terraformOptions, "domain_name"), "domain_name")
			assert.Equal(t, testCase.localPart, terraform.Output(t, terraformOptions, "local_part"), "local_part")
			assert.Equal(t, testCase.want.Address, terraform.Output(t, terraformOptions, "address"), "address")
			assert.Equal(t, testCase.want.ExpiresOn, terraform.Output(t, terraformOptions, "expires_on"), "expires_on")
			assert.Equal(t, strconv.FormatBool(testCase.want.Expirable), terraform.Output(t, terraformOptions, "expirable"), "expirable")
		})
	}
}

func TestMailboxResource_Using_RecoveryEmail(t *testing.T) {
	testCases := map[string]struct {
		localPart string
		domain    string
		state     []model.Mailbox
		want      model.Mailbox
	}{
		"single": {
			localPart: "some",
			domain:    "example.com",
			state:     []model.Mailbox{},
			want: model.Mailbox{
				LocalPart:             "some",
				DomainName:            "example.com",
				Address:               "some@example.com",
				IsInternal:            true,
				Expirable:             false,
				ExpiresOn:             "",
				RemoveUponExpiry:      false,
				PasswordRecoveryEmail: "someone@example.com",
			},
		},
		"multiple": {
			localPart: "some",
			domain:    "example.com",
			state: []model.Mailbox{
				{
					LocalPart:        "other",
					DomainName:       "example.com",
					Address:          "other@example.com",
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: model.Mailbox{
				LocalPart:             "some",
				DomainName:            "example.com",
				Address:               "some@example.com",
				IsInternal:            true,
				Expirable:             false,
				ExpiresOn:             "",
				RemoveUponExpiry:      false,
				PasswordRecoveryEmail: "someone@example.com",
			},
		},
		"idna": {
			localPart: "test",
			domain:    "hoß.de",
			state:     []model.Mailbox{},
			want: model.Mailbox{
				LocalPart:             "test",
				DomainName:            "xn--ho-hia.de",
				Address:               "test@xn--ho-hia.de",
				Expirable:             false,
				ExpiresOn:             "",
				RemoveUponExpiry:      false,
				PasswordRecoveryEmail: "someone@example.com",
			},
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state}))
			defer server.Close()

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				TerraformDir: "../resources/migadu_mailbox/invitation",
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
			assert.Equal(t, tt.want.PasswordRecoveryEmail, terraform.Output(t, terraformOptions, "password_recovery_email"), "password_recovery_email")
			assert.Equal(t, tt.want.ExpiresOn, terraform.Output(t, terraformOptions, "expires_on"), "expires_on")
			assert.Equal(t, strconv.FormatBool(tt.want.Expirable), terraform.Output(t, terraformOptions, "expirable"), "expirable")
		})
	}
}
