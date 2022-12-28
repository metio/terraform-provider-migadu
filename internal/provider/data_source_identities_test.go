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

func TestIdentitiesDataSource_Read(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		want       *client.Identities
		error      string
	}{
		{
			name:       "empty",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusOK,
			want: &client.Identities{
				Identities: []client.Identity{},
			},
		},
		{
			name:       "single",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusOK,
			want: &client.Identities{
				Identities: []client.Identity{
					{
						LocalPart:            "test",
						DomainName:           "example.com",
						Address:              "test@example.com",
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
			},
		},
		{
			name:       "multiple",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusOK,
			want: &client.Identities{
				Identities: []client.Identity{
					{
						LocalPart:            "test",
						DomainName:           "example.com",
						Address:              "test@example.com",
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
					{
						LocalPart:  "other",
						DomainName: "example.com",
						Address:    "other@example.com",
						Name:       "Other Identity",
					},
				},
			},
		},
		{
			name:       "idna",
			domain:     "hoß.de",
			localPart:  "test",
			statusCode: http.StatusOK,
			want: &client.Identities{
				Identities: []client.Identity{
					{
						LocalPart:  "test",
						DomainName: "xn--ho-hia.de",
						Address:    "test@xn--ho-hia.de",
					},
				},
			},
		},
		{
			name:       "error-401",
			domain:     "example.com",
			localPart:  "not-found",
			statusCode: http.StatusUnauthorized,
			want:       nil,
			error:      "Request failed with: status: 401",
		},
		{
			name:       "error-404",
			domain:     "example.com",
			localPart:  "not-found",
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
					data "migadu_identities" "test" {
						domain_name = "%s"
						local_part  = "%s"
					}
				`, tt.domain, tt.localPart)

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
								resource.TestCheckResourceAttr("data.migadu_identities.test", "domain_name", tt.domain),
								resource.TestCheckResourceAttr("data.migadu_identities.test", "local_part", tt.localPart),
								resource.TestCheckResourceAttr("data.migadu_identities.test", "identities.#", fmt.Sprintf("%v", len(tt.want.Identities))),
								resource.TestCheckResourceAttr("data.migadu_identities.test", "id", fmt.Sprintf("%s@%s", tt.localPart, tt.domain)),
							),
						},
					},
				})
			}
		})
	}
}
