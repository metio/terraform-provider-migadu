/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-migadu/internal/client"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestMailboxesDataSource_Read(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		want       *client.Mailboxes
		error      string
	}{
		{
			name:       "empty",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want: &client.Mailboxes{
				Mailboxes: []client.Mailbox{},
			},
		},
		{
			name:       "single",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want: &client.Mailboxes{
				Mailboxes: []client.Mailbox{
					{
						LocalPart:             "test",
						DomainName:            "example.com",
						Address:               "test@example.com",
						Name:                  "Some Name",
						IsInternal:            false,
						MaySend:               true,
						MayReceive:            false,
						MayAccessImap:         true,
						MayAccessPop3:         false,
						MayAccessManageSieve:  true,
						PasswordMethod:        "",
						Password:              "hunter2",
						PasswordRecoveryEmail: "recovery@example.com",
						SpamAction:            "",
						SpamAggressiveness:    "",
						Expirable:             true,
						ExpiresOn:             "",
						RemoveUponExpiry:      false,
						SenderDenyList:        []string{},
						SenderAllowList:       []string{},
						RecipientDenyList:     []string{},
						AutoRespondActive:     true,
						AutoRespondSubject:    "kthxbye",
						AutoRespondBody:       "",
						AutoRespondExpiresOn:  "",
						FooterActive:          false,
						FooterPlainBody:       "",
						FooterHtmlBody:        "",
						StorageUsage:          0.5,
						Delegations:           []string{},
						Identities:            []string{},
					},
				},
			},
		},
		{
			name:       "multiple",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want: &client.Mailboxes{
				Mailboxes: []client.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "example.com",
						Address:    "test@example.com",
						Name:       "Some Name",
					},
					{
						LocalPart:  "other",
						DomainName: "example.com",
						Address:    "other@example.com",
						Name:       "Other Name",
					},
				},
			},
		},
		{
			name:       "idna",
			domain:     "hoß.de",
			statusCode: http.StatusOK,
			want: &client.Mailboxes{
				Mailboxes: []client.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "xn--ho-hia.de",
						Address:    "test@xn--ho-hia.de",
						Name:       "Some Name",
					},
				},
			},
		},
		{
			name:       "error-401",
			domain:     "example.com",
			statusCode: http.StatusUnauthorized,
			want:       nil,
			error:      "Request failed with: status: 401",
		},
		{
			name:       "error-404",
			domain:     "example.com",
			statusCode: http.StatusNotFound,
			want:       nil,
			error:      "Request failed with: status: 404",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				bytes, err := json.Marshal(tt.want)
				if err != nil {
					t.Errorf("Could not serialize data")
				}
				w.Write(bytes)
			}))
			defer server.Close()

			config := providerConfig(server.URL) + fmt.Sprintf(`
					data "migadu_mailboxes" "test" {
						domain_name = "%s"
					}
				`, tt.domain)

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
								resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "domain_name", tt.domain),
								resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "mailboxes.#", fmt.Sprintf("%v", len(tt.want.Mailboxes))),
								resource.TestCheckResourceAttr("data.migadu_mailboxes.test", "id", fmt.Sprintf("%s", tt.domain)),
							),
						},
					},
				})
			}
		})
	}
}
