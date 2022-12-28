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

func TestMigaduClient_GetRewrites(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		wantedDomain string
		statusCode   int
		want         *Rewrites
		wantErr      bool
	}{
		{
			name:         "empty",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want:         &Rewrites{},
			wantErr:      false,
		},
		{
			name:         "single",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Rewrites{
				[]Rewrite{
					{
						DomainName:    "example.com",
						Name:          "test",
						LocalPartRule: "rule-*",
						OrderNum:      1,
						Destinations: []string{
							"another@example.com",
						},
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
			want: &Rewrites{
				[]Rewrite{
					{
						DomainName:    "example.com",
						Name:          "test",
						LocalPartRule: "rule-*",
						OrderNum:      1,
						Destinations: []string{
							"another@example.com",
						},
					},
					{
						DomainName:    "example.com",
						Name:          "another",
						LocalPartRule: "rule*",
						OrderNum:      3,
						Destinations: []string{
							"some@example.com",
							"other@example.com",
						},
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
			want:         &Rewrites{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/domains/%s/rewrites", tt.wantedDomain) {
					t.Errorf("Expected to request '/domains/%s/rewrites', got: %s", tt.wantedDomain, r.URL.Path)
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

			got, err := c.GetRewrites(tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRewrites() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRewrites() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_GetRewrite(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		slug         string
		wantedDomain string
		statusCode   int
		want         *Rewrite
		wantErr      bool
	}{
		{
			name:         "single",
			domain:       "example.com",
			slug:         "slug",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Rewrite{
				DomainName:    "example.com",
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"securitu@example.com",
				},
			},
			wantErr: false,
		},
		{
			name:         "idna",
			domain:       "hoß.de",
			slug:         "slug",
			wantedDomain: "xn--ho-hia.de",
			statusCode:   http.StatusOK,
			want: &Rewrite{
				DomainName:    "xn--ho-hia.de",
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"securitu@xn--ho-hia.de",
				},
			},
			wantErr: false,
		},
		{
			name:         "error-404",
			domain:       "example.com",
			slug:         "slug",
			wantedDomain: "example.com",
			statusCode:   http.StatusNotFound,
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/domains/%s/rewrites/%s", tt.wantedDomain, tt.slug) {
					t.Errorf("Expected to request '/domains/%s/rewrites/%s', got: %s", tt.wantedDomain, tt.slug, r.URL.Path)
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

			got, err := c.GetRewrite(tt.domain, tt.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRewrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRewrite() got = %v, want %v", got, tt.want)
			}
		})
	}
}
