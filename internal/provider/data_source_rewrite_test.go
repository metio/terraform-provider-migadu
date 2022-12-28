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

func TestRewriteDataSource_Read(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		slug       string
		statusCode int
		want       *client.Rewrite
		error      string
	}{
		{
			name:       "empty",
			domain:     "example.com",
			slug:       "test",
			statusCode: http.StatusOK,
			want:       &client.Rewrite{},
		},
		{
			name:       "single",
			domain:     "example.com",
			slug:       "test",
			statusCode: http.StatusOK,
			want: &client.Rewrite{
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
			name:       "idna",
			domain:     "hoß.de",
			slug:       "test",
			statusCode: http.StatusOK,
			want: &client.Rewrite{
				DomainName:    "xn--ho-hia.de",
				Name:          "test",
				LocalPartRule: "prefix-*",
				OrderNum:      0,
				Destinations: []string{
					"dest@xn--ho-hia.de",
				},
			},
		},
		{
			name:       "error-401",
			domain:     "example.com",
			slug:       "test",
			statusCode: http.StatusUnauthorized,
			want:       nil,
			error:      "Request failed with: status: 401",
		},
		{
			name:       "error-404",
			domain:     "example.com",
			slug:       "test",
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
					data "migadu_rewrite" "test" {
						domain_name = "%s"
						name        = "%s"
					}
				`, tt.domain, tt.slug)

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
								resource.TestCheckResourceAttr("data.migadu_rewrite.test", "domain_name", tt.domain),
								resource.TestCheckResourceAttr("data.migadu_rewrite.test", "name", tt.slug),
								resource.TestCheckResourceAttr("data.migadu_rewrite.test", "id", fmt.Sprintf("%s@%s", tt.slug, tt.domain)),
							),
						},
					},
				})
			}
		})
	}
}
