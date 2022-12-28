/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMigaduClient_GetAliases(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		wantedDomain string
		statusCode   int
		want         *Aliases
		wantErr      bool
	}{
		{
			name:         "empty",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want:         &Aliases{Aliases: []Alias{}},
			wantErr:      false,
		},
		{
			name:         "single",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Aliases{
				Aliases: []Alias{
					{
						LocalPart:        "local",
						DomainName:       "example.com",
						Address:          "another",
						Destinations:     []string{},
						IsInternal:       false,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "multiple",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Aliases{
				Aliases: []Alias{
					{
						LocalPart:        "local",
						DomainName:       "example.com",
						Address:          "another",
						Destinations:     []string{},
						IsInternal:       false,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: true,
					},
					{
						LocalPart:  "test",
						DomainName: "example.com",
						Address:    "address",
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
			wantErr: false,
		},
		{
			name:         "error-404",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusNotFound,
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "idna",
			domain:       "hoß.de",
			wantedDomain: "xn--ho-hia.de",
			statusCode:   http.StatusOK,
			want:         &Aliases{Aliases: []Alias{}},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/domains/%s/aliases", tt.wantedDomain) {
					t.Errorf("Expected to request '/domains/%s/aliases', got: %s", tt.wantedDomain, r.URL.Path)
				}
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

			c := newTestClient(server.URL)

			got, err := c.GetAliases(tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAliases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAliases() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_GetAlias(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		localPart    string
		wantedDomain string
		statusCode   int
		want         *Alias
		wantErr      bool
	}{
		{
			name:         "empty",
			domain:       "example.com",
			localPart:    "test",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want:         &Alias{},
			wantErr:      false,
		},
		{
			name:         "single",
			domain:       "example.com",
			localPart:    "test",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Alias{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Destinations: []string{
					"different@example.com",
				},
				IsInternal:       false,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
			wantErr: false,
		},
		{
			name:         "multiple",
			domain:       "example.com",
			localPart:    "test",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Alias{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Destinations: []string{
					"different@example.com",
					"another@example.com",
				},
			},
			wantErr: false,
		},
		{
			name:         "error-404",
			domain:       "example.com",
			localPart:    "test",
			wantedDomain: "example.com",
			statusCode:   http.StatusNotFound,
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "idna",
			domain:       "hoß.de",
			localPart:    "seb",
			wantedDomain: "xn--ho-hia.de",
			statusCode:   http.StatusOK,
			want:         &Alias{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/domains/%s/aliases/%s", tt.wantedDomain, tt.localPart) {
					t.Errorf("Expected to request '/domains/%s/aliases/%s', got: %s", tt.wantedDomain, tt.localPart, r.URL.Path)
				}
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

			c := newTestClient(server.URL)

			got, err := c.GetAlias(tt.domain, tt.localPart)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlias() got = %v, want %v", got, tt.want)
			}
		})
	}
}
