/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-migadu/internal/migadu/model"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestAliasesDataSource_Read(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		want       *model.Aliases
		error      string
	}{
		{
			name:       "empty",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want:       &model.Aliases{Aliases: []model.Alias{}},
		},
		{
			name:       "single",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want: &model.Aliases{
				Aliases: []model.Alias{
					{
						LocalPart:        "local",
						DomainName:       "example.com",
						Address:          "local@example.com",
						Destinations:     []string{},
						IsInternal:       false,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: true,
					},
				},
			},
		},
		{
			name:       "multiple",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want: &model.Aliases{
				Aliases: []model.Alias{
					{
						LocalPart:        "local",
						DomainName:       "example.com",
						Address:          "local@example.com",
						Destinations:     []string{},
						IsInternal:       false,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: true,
					},
					{
						LocalPart:  "test",
						DomainName: "example.com",
						Address:    "test@example.com",
						Destinations: []string{
							"destination@example.com",
						},
						IsInternal:       true,
						Expirable:        true,
						ExpiresOn:        "",
						RemoveUponExpiry: false,
					},
				},
			},
		},
		{
			name:       "idna",
			domain:     "hoß.de",
			statusCode: http.StatusOK,
			want: &model.Aliases{
				Aliases: []model.Alias{
					{
						LocalPart:  "another",
						DomainName: "xn--ho-hia.de",
						Address:    "another@xn--ho-hia.de",
					},
				},
			},
		},
		{
			name:       "error-401",
			domain:     "example.com",
			statusCode: http.StatusUnauthorized,
			want:       nil,
			error:      "status: 401",
		},
		{
			name:       "error-404",
			domain:     "example.com",
			statusCode: http.StatusNotFound,
			want:       nil,
			error:      "status: 404",
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
				_, err = w.Write(bytes)
				if err != nil {
					t.Errorf("Could not write data")
				}
			}))
			defer server.Close()

			config := providerConfig(server.URL) + fmt.Sprintf(`
					data "migadu_aliases" "test" {
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
								resource.TestCheckResourceAttr("data.migadu_aliases.test", "domain_name", tt.domain),
								resource.TestCheckResourceAttr("data.migadu_aliases.test", "address_aliases.#", fmt.Sprintf("%v", len(tt.want.Aliases))),
								resource.TestCheckResourceAttr("data.migadu_aliases.test", "id", tt.domain),
							),
						},
					},
				})
			}
		})
	}
}
