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

func TestRewritesDataSource_Read(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		want       *client.Rewrites
		error      string
	}{
		{
			name:       "empty",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want: &client.Rewrites{
				Rewrites: []client.Rewrite{},
			},
		},
		{
			name:       "single",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want: &client.Rewrites{
				Rewrites: []client.Rewrite{
					{
						DomainName:    "example.com",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest@example.com",
						},
					},
				},
			},
		},
		{
			name:       "multiple",
			domain:     "example.com",
			statusCode: http.StatusOK,
			want: &client.Rewrites{
				Rewrites: []client.Rewrite{
					{
						DomainName:    "example.com",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest1@example.com",
							"dest2@example.com",
						},
					},
					{
						DomainName:    "example.com",
						Name:          "Other Rule",
						LocalPartRule: "*-suffix",
						OrderNum:      1,
						Destinations: []string{
							"dest3@example.com",
							"dest4@example.com",
						},
					},
				},
			},
		},
		{
			name:       "idna",
			domain:     "hoß.de",
			statusCode: http.StatusOK,
			want: &client.Rewrites{
				Rewrites: []client.Rewrite{
					{
						DomainName:    "xn--ho-hia.de",
						Name:          "Some Rule",
						LocalPartRule: "prefix-*",
						OrderNum:      0,
						Destinations: []string{
							"dest@xn--ho-hia.de",
						},
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
					data "migadu_rewrites" "test" {
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
								resource.TestCheckResourceAttr("data.migadu_rewrites.test", "domain_name", tt.domain),
								resource.TestCheckResourceAttr("data.migadu_rewrites.test", "rewrites.#", fmt.Sprintf("%v", len(tt.want.Rewrites))),
								resource.TestCheckResourceAttr("data.migadu_rewrites.test", "id", fmt.Sprintf("%s", tt.domain)),
							),
						},
					},
				})
			}
		})
	}
}
