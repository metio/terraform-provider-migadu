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
	"golang.org/x/net/idna"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestAliasResource(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		localPart    string
		destinations []string
		want         *client.Alias
		error        string
	}{
		{
			name:      "no-domain",
			domain:    "",
			localPart: "test",
			destinations: []string{
				"other@example.com",
			},
			want:  &client.Alias{},
			error: "Attribute domain_name string length must be at least 1",
		},
		{
			name:      "no-local-part",
			domain:    "example.com",
			localPart: "",
			destinations: []string{
				"other@example.com",
			},
			want:  &client.Alias{},
			error: "Attribute local_part string length must be at least 1",
		},
		{
			name:         "no-destinations",
			domain:       "example.com",
			localPart:    "test",
			destinations: []string{},
			want:         &client.Alias{},
			error:        "Attribute destinations list must contain at least 1 elements",
		},
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			destinations: []string{
				"other@example.com",
			},
			want: &client.Alias{
				Address: "test@example.com",
				Destinations: []string{
					"other@example.com",
				},
				IsInternal:       false,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			destinations: []string{
				"other@example.com",
				"some@example.com",
			},
			want: &client.Alias{
				Address: "test@example.com",
				Destinations: []string{
					"other@example.com",
					"some@example.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(aliasSimulator(t))
			defer server.Close()

			config := providerConfig(server.URL) + fmt.Sprintf(`
					resource "migadu_alias" "test" {
						domain_name  = "%s"
						local_part   = "%s"
						destinations = %s
					}
				`, tt.domain, tt.localPart, strings.ReplaceAll(fmt.Sprintf("%+q", tt.destinations), "\" \"", "\",\""))

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
								resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", tt.domain),
								resource.TestCheckResourceAttr("migadu_alias.test", "local_part", tt.localPart),
								resource.TestCheckResourceAttr("migadu_alias.test", "address", tt.want.Address),
								resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", fmt.Sprintf("%v", len(tt.want.Destinations))),
								resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", tt.want.Destinations[0]),
								resource.TestCheckResourceAttr("migadu_alias.test", "id", fmt.Sprintf("%s@%s", tt.localPart, tt.domain)),
							),
						},
						{
							ResourceName:      "migadu_alias.test",
							ImportState:       true,
							ImportStateVerify: true,
						},
					},
				})
			}
		})
	}
}

func TestAliasResource_IDN_ASCII(t *testing.T) {
	server := httptest.NewServer(aliasSimulator(t))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(server.URL) + `
					resource "migadu_alias" "test" {
						domain_name  = "hoß.de"
						local_part   = "test"
						destinations = ["first@xn--ho-hia.de", "second@xn--ho-hia.de"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", "hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "local_part", "test"),
					resource.TestCheckResourceAttr("migadu_alias.test", "address", "test@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", "2"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", "first@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.1", "second@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.#", "2"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.0", "first@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.1", "second@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "id", "test@hoß.de"),
				),
			},
			{
				ResourceName:      "migadu_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig(server.URL) + `
					resource "migadu_alias" "test" {
						domain_name      = "hoß.de"
						local_part       = "test"
						destinations_idn = ["third@xn--ho-hia.de"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", "hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "local_part", "test"),
					resource.TestCheckResourceAttr("migadu_alias.test", "address", "test@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", "1"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", "third@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.#", "1"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.0", "third@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "id", "test@hoß.de"),
				),
			},
		},
	})
}

func TestAliasResource_IDN_Unicode(t *testing.T) {
	server := httptest.NewServer(aliasSimulator(t))
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(server.URL) + `
					resource "migadu_alias" "test" {
						domain_name      = "hoß.de"
						local_part       = "test"
						destinations_idn = ["first@hoß.de", "second@hoß.de"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", "hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "local_part", "test"),
					resource.TestCheckResourceAttr("migadu_alias.test", "address", "test@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", "2"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", "first@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.1", "second@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.#", "2"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.0", "first@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.1", "second@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "id", "test@hoß.de"),
				),
			},
			{
				ResourceName:      "migadu_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig(server.URL) + `
					resource "migadu_alias" "test" {
						domain_name      = "hoß.de"
						local_part       = "test"
						destinations_idn = ["third@hoß.de"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("migadu_alias.test", "domain_name", "hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "local_part", "test"),
					resource.TestCheckResourceAttr("migadu_alias.test", "address", "test@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.#", "1"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations.0", "third@xn--ho-hia.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.#", "1"),
					resource.TestCheckResourceAttr("migadu_alias.test", "destinations_idn.0", "third@hoß.de"),
					resource.TestCheckResourceAttr("migadu_alias.test", "id", "test@hoß.de"),
				),
			},
		},
	})
}

func aliasSimulator(t *testing.T) http.HandlerFunc {
	var aliases []client.Alias
	urlPattern := regexp.MustCompile("/domains/(.*)/aliases/?(.*)?")

	return func(w http.ResponseWriter, r *http.Request) {
		matches := urlPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			t.Errorf("Expected to request to match %s, got: %s", urlPattern, r.URL.Path)
		}
		domain := matches[1]
		localPart := matches[2]

		if r.Method == http.MethodPost {
			handleCreateAlias(w, r, t, &aliases, domain)
		}
		if r.Method == http.MethodPut {
			handleUpdateAlias(w, r, t, &aliases, domain, localPart)
		}
		if r.Method == http.MethodDelete {
			handleDeleteAlias(w, r, t, &aliases, domain, localPart)
		}
		if r.Method == http.MethodGet {
			handleGetAlias(w, r, t, &aliases, domain, localPart)
		}
	}
}

func handleGetAlias(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]client.Alias, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/aliases/%s', got: %s", domain, localPart, r.URL.Path)
	}

	missing := true
	for _, alias := range *aliases {
		if alias.DomainName == domain && alias.LocalPart == localPart {
			missing = false
			w.WriteHeader(http.StatusOK)
			bytes, err := json.Marshal(alias)
			if err != nil {
				t.Errorf("Could not marshall alias")
			}
			_, err = w.Write(bytes)
			if err != nil {
				t.Errorf("Could not write data")
			}
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleDeleteAlias(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]client.Alias, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/aliases/%s', got: %s", domain, localPart, r.URL.Path)
	}

	missing := true
	for index, alias := range *aliases {
		if alias.DomainName == domain && alias.LocalPart == localPart {
			missing = false
			w.WriteHeader(http.StatusOK)
			c := *aliases
			c[index] = c[len(c)-1]
			*aliases = c[:len(c)-1]

			bytes, err := json.Marshal(alias)
			if err != nil {
				t.Errorf("Could not marshall alias")
			}
			_, err = w.Write(bytes)
			if err != nil {
				t.Errorf("Could not write data")
			}
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleUpdateAlias(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]client.Alias, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/aliases/%s', got: %s", domain, localPart, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	requestAlias := client.Alias{}
	err = json.Unmarshal(requestBody, &requestAlias)
	if err != nil {
		t.Errorf("Could not unmarshall alias")
	}

	requestAlias.DomainName = domain
	requestAlias.Address = fmt.Sprintf("%s@%s", requestAlias.LocalPart, domain)

	var asciiDestinations []string
	for _, dest := range requestAlias.Destinations {
		parts := strings.Split(dest, "@")

		ascii, err := idna.ToASCII(parts[1])
		if err != nil {
			t.Errorf("could not convert to punycode")
		}
		asciiDestinations = append(asciiDestinations, fmt.Sprintf("%s@%s", parts[0], ascii))
	}
	requestAlias.Destinations = asciiDestinations

	missing := true
	for index, alias := range *aliases {
		if alias.DomainName == domain && alias.LocalPart == localPart {
			missing = false
			w.WriteHeader(http.StatusOK)
			c := *aliases
			c[index] = requestAlias
			*aliases = c

			bytes, err := json.Marshal(requestAlias)
			if err != nil {
				t.Errorf("Could not marshall alias")
			}
			_, err = w.Write(bytes)
			if err != nil {
				t.Errorf("Could not write data")
			}
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleCreateAlias(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]client.Alias, domain string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases", domain) {
		t.Errorf("Expected to request '/domains/%s/aliases', got: %s", domain, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	alias := client.Alias{}
	err = json.Unmarshal(requestBody, &alias)
	if err != nil {
		t.Errorf("Could not unmarshall alias")
	}
	alias.DomainName = domain
	alias.Address = fmt.Sprintf("%s@%s", alias.LocalPart, domain)

	var asciiDestinations []string
	for _, dest := range alias.Destinations {
		parts := strings.Split(dest, "@")

		ascii, err := idna.ToASCII(parts[1])
		if err != nil {
			t.Errorf("could not convert to punycode")
		}
		asciiDestinations = append(asciiDestinations, fmt.Sprintf("%s@%s", parts[0], ascii))
	}
	alias.Destinations = asciiDestinations

	*aliases = append(*aliases, alias)

	responseBody, err := json.Marshal(alias)
	_, err = w.Write(responseBody)
	if err != nil {
		t.Errorf("Could not write data")
	}
}
