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

func TestIdentityDataSource_Read(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		identity   string
		statusCode int
		want       *client.Identity
		error      string
	}{
		{
			name:       "empty",
			domain:     "example.com",
			localPart:  "test",
			identity:   "someone",
			statusCode: http.StatusOK,
			want:       &client.Identity{},
		},
		{
			name:       "single",
			domain:     "example.com",
			localPart:  "test",
			identity:   "someone",
			statusCode: http.StatusOK,
			want: &client.Identity{
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
		{
			name:       "idna",
			domain:     "hoß.de",
			localPart:  "test",
			identity:   "someone",
			statusCode: http.StatusOK,
			want: &client.Identity{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
			},
		},
		{
			name:       "error-401",
			domain:     "example.com",
			localPart:  "not-found",
			identity:   "someone",
			statusCode: http.StatusUnauthorized,
			want:       nil,
			error:      "Request failed with: status: 401",
		},
		{
			name:       "error-404",
			domain:     "example.com",
			localPart:  "not-found",
			identity:   "someone",
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
					data "migadu_identity" "test" {
						domain_name = "%s"
						local_part  = "%s"
						identity    = "%s"
					}
				`, tt.domain, tt.localPart, tt.identity)

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
								resource.TestCheckResourceAttr("data.migadu_identity.test", "domain_name", tt.domain),
								resource.TestCheckResourceAttr("data.migadu_identity.test", "local_part", tt.localPart),
								resource.TestCheckResourceAttr("data.migadu_identity.test", "identity", tt.identity),
								resource.TestCheckResourceAttr("data.migadu_identity.test", "id", fmt.Sprintf("%s@%s/%s", tt.localPart, tt.domain, tt.identity)),
							),
						},
					},
				})
			}
		})
	}
}
