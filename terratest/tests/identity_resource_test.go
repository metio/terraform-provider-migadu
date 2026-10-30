/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package acceptance_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/metio/migadu-client.go/model"
	"github.com/metio/migadu-client.go/simulator"
	"github.com/stretchr/testify/assert"
)

func TestIdentityResource(t *testing.T) {
	testCases := map[string]struct {
		domain    string
		localPart string
		identity  string
		state     []model.Identity
		want      model.Identity
	}{
		"single": {
			domain:    "example.com",
			localPart: "some",
			identity:  "other",
			state:     []model.Identity{},
			want: model.Identity{
				LocalPart:  "some",
				DomainName: "example.com",
				Address:    "other@example.com",
			},
		},
		"multiple": {
			domain:    "example.com",
			localPart: "some",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "different",
					DomainName: "example.com",
					Address:    "other@example.com",
				},
			},
			want: model.Identity{
				LocalPart:  "some",
				DomainName: "example.com",
				Address:    "other@example.com",
			},
		},
		"idna": {
			domain:    "hoß.de",
			localPart: "test",
			identity:  "other",
			state:     []model.Identity{},
			want: model.Identity{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "other@xn--ho-hia.de",
			},
		},
	}
	for name, testCase := range testCases {
		for _, passwordUse := range []string{"custom", "mailbox", "none"} {
			t.Run(fmt.Sprintf("%s/%s", passwordUse, name), func(t *testing.T) {
				server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: testCase.state}))
				defer server.Close()

				terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
					TerraformDir: fmt.Sprintf("../resources/migadu_identity/%s", passwordUse),
					Vars: map[string]interface{}{
						"endpoint":    server.URL,
						"domain_name": testCase.domain,
						"local_part":  testCase.localPart,
						"identity":    testCase.identity,
					},
				})

				defer terraform.DestroyContext(t, t.Context(), terraformOptions)
				terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

				assert.Equal(t, fmt.Sprintf("%s@%s/%s", testCase.localPart, testCase.domain, testCase.identity), terraform.OutputContext(t, t.Context(), terraformOptions, "id"), "id")
				assert.Equal(t, testCase.domain, terraform.OutputContext(t, t.Context(), terraformOptions, "domain_name"), "domain_name")
				assert.Equal(t, testCase.localPart, terraform.OutputContext(t, t.Context(), terraformOptions, "local_part"), "local_part")
				assert.Equal(t, testCase.want.Address, terraform.OutputContext(t, t.Context(), terraformOptions, "address"), "address")
			})
		}
	}
}
