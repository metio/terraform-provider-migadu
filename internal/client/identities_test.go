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

func TestMigaduClient_GetIdentities(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		localPart    string
		wantedDomain string
		statusCode   int
		want         *Identities
		wantErr      bool
	}{
		{
			name:         "empty",
			domain:       "example.com",
			localPart:    "test",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want:         &Identities{},
			wantErr:      false,
		},
		{
			name:         "single",
			domain:       "example.com",
			localPart:    "test",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Identities{
				[]Identity{
					{
						LocalPart:            "test",
						DomainName:           "example.com",
						Address:              "test@example.com",
						Name:                 "Some Name",
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
			wantErr: false,
		},
		{
			name:         "multiple",
			domain:       "example.com",
			localPart:    "test",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Identities{
				[]Identity{
					{
						LocalPart:            "test",
						DomainName:           "example.com",
						Address:              "test@example.com",
						Name:                 "Some Name",
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
						LocalPart:            "some",
						DomainName:           "example.com",
						Address:              "some@example.com",
						Name:                 "Someone Else",
						MaySend:              false,
						MayReceive:           false,
						MayAccessImap:        false,
						MayAccessPop3:        false,
						MayAccessManageSieve: false,
						Password:             "hunter2",
						FooterActive:         true,
						FooterPlainBody:      "this is my footer",
						FooterHtmlBody:       "<strong>have a nice day</strong>",
					},
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
			localPart:    "test",
			wantedDomain: "xn--ho-hia.de",
			statusCode:   http.StatusOK,
			want:         &Identities{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s/identities", tt.wantedDomain, tt.localPart) {
					t.Errorf("Expected to request '/domains/%s/mailboxes/%s/identities', got: %s", tt.wantedDomain, tt.localPart, r.URL.Path)
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

			got, err := c.GetIdentities(tt.domain, tt.localPart)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIdentities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIdentities() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_GetIdentity(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		localPart    string
		identity     string
		wantedDomain string
		statusCode   int
		want         *Identity
		wantErr      bool
	}{
		{
			name:         "single",
			domain:       "example.com",
			localPart:    "test",
			identity:     "other",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Identity{
				LocalPart:            "other",
				DomainName:           "example.com",
				Address:              "other@example.com",
				Name:                 "Some Name",
				MaySend:              false,
				MayReceive:           false,
				MayAccessImap:        false,
				MayAccessPop3:        false,
				MayAccessManageSieve: false,
				Password:             "",
				FooterActive:         false,
				FooterPlainBody:      "",
				FooterHtmlBody:       "",
			},
			wantErr: false,
		},
		{
			name:         "error-404",
			domain:       "example.com",
			localPart:    "test",
			identity:     "other",
			wantedDomain: "example.com",
			statusCode:   http.StatusNotFound,
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "idna",
			domain:       "hoß.de",
			localPart:    "test",
			identity:     "other",
			wantedDomain: "xn--ho-hia.de",
			statusCode:   http.StatusOK,
			want: &Identity{
				LocalPart:  "other",
				DomainName: "xn--ho-hia.de",
				Address:    "other@xn--ho-hia.de",
				Name:       "Some Name",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s/identities/%s", tt.wantedDomain, tt.localPart, tt.identity) {
					t.Errorf("Expected to request '/domains/%s/mailboxes/%s/identities/%s', got: %s", tt.wantedDomain, tt.localPart, tt.identity, r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				bytes, err := json.Marshal(tt.want)
				if err != nil {
					t.Errorf("Could not serialize data")
				}
				w.Write(bytes)
			}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetIdentity(tt.domain, tt.localPart, tt.identity)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIdentity() got = %v, want %v", got, tt.want)
			}
		})
	}
}
