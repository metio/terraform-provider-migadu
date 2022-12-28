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

func TestMigaduClient_GetMailboxes(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		wantedDomain string
		statusCode   int
		want         *Mailboxes
		wantErr      bool
	}{
		{
			name:         "empty",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want:         &Mailboxes{},
			wantErr:      false,
		},
		{
			name:         "single",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Mailboxes{
				[]Mailbox{
					{
						LocalPart:             "test",
						DomainName:            "example.com",
						Address:               "test@example.com",
						Name:                  "test",
						IsInternal:            false,
						MaySend:               true,
						MayReceive:            true,
						MayAccessImap:         true,
						MayAccessPop3:         true,
						MayAccessManageSieve:  true,
						PasswordMethod:        "",
						Password:              "secret",
						PasswordRecoveryEmail: "",
						SpamAction:            "",
						SpamAggressiveness:    "",
						Expirable:             false,
						ExpiresOn:             "",
						RemoveUponExpiry:      false,
						SenderDenyList:        []string{},
						SenderAllowList:       []string{},
						RecipientDenyList:     []string{},
						AutoRespondActive:     false,
						AutoRespondSubject:    "",
						AutoRespondBody:       "",
						AutoRespondExpiresOn:  "",
						FooterActive:          false,
						FooterPlainBody:       "",
						FooterHtmlBody:        "",
						StorageUsage:          0.0,
						Delegations:           []string{},
						Identities:            []string{}},
				},
			},
			wantErr: false,
		},
		{
			name:         "multiple",
			domain:       "example.com",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Mailboxes{
				[]Mailbox{
					{
						LocalPart:             "test",
						DomainName:            "example.com",
						Address:               "test@example.com",
						Name:                  "test",
						IsInternal:            false,
						MaySend:               true,
						MayReceive:            true,
						MayAccessImap:         true,
						MayAccessPop3:         true,
						MayAccessManageSieve:  true,
						PasswordMethod:        "",
						Password:              "secret",
						PasswordRecoveryEmail: "",
						SpamAction:            "",
						SpamAggressiveness:    "",
						Expirable:             false,
						ExpiresOn:             "",
						RemoveUponExpiry:      false,
						SenderDenyList:        []string{},
						SenderAllowList:       []string{},
						RecipientDenyList:     []string{},
						AutoRespondActive:     false,
						AutoRespondSubject:    "",
						AutoRespondBody:       "",
						AutoRespondExpiresOn:  "",
						FooterActive:          false,
						FooterPlainBody:       "",
						FooterHtmlBody:        "",
						StorageUsage:          0.0,
						Delegations:           []string{},
						Identities:            []string{},
					},
					{
						LocalPart:  "another",
						DomainName: "example.com",
						Address:    "another@example.com",
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
			want:         &Mailboxes{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes", tt.wantedDomain) {
					t.Errorf("Expected to request '/domains/%s/mailboxes', got: %s", tt.wantedDomain, r.URL.Path)
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

			got, err := c.GetMailboxes(tt.domain)
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

func TestMigaduClient_GetMailbox(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		localPart    string
		wantedDomain string
		statusCode   int
		want         *Mailbox
		wantErr      bool
	}{
		{
			name:         "single",
			domain:       "example.com",
			localPart:    "test",
			wantedDomain: "example.com",
			statusCode:   http.StatusOK,
			want: &Mailbox{
				LocalPart:             "test",
				DomainName:            "example.com",
				Address:               "test@example.com",
				Name:                  "Some Name",
				IsInternal:            false,
				MaySend:               false,
				MayReceive:            false,
				MayAccessImap:         false,
				MayAccessPop3:         false,
				MayAccessManageSieve:  false,
				PasswordMethod:        "",
				Password:              "",
				PasswordRecoveryEmail: "",
				SpamAction:            "",
				SpamAggressiveness:    "",
				Expirable:             false,
				ExpiresOn:             "",
				RemoveUponExpiry:      false,
				SenderDenyList:        nil,
				SenderAllowList:       nil,
				RecipientDenyList:     nil,
				AutoRespondActive:     false,
				AutoRespondSubject:    "",
				AutoRespondBody:       "",
				AutoRespondExpiresOn:  "",
				FooterActive:          false,
				FooterPlainBody:       "",
				FooterHtmlBody:        "",
				StorageUsage:          0,
				Delegations:           nil,
				Identities:            nil,
			},
			wantErr: false,
		},
		{
			name:         "idna",
			domain:       "hoß.de",
			localPart:    "test",
			wantedDomain: "xn--ho-hia.de",
			statusCode:   http.StatusOK,
			want: &Mailbox{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Name:       "Some Name",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s", tt.wantedDomain, tt.localPart) {
					t.Errorf("Expected to request '/domains/%s/mailboxes/%s', got: %s", tt.wantedDomain, tt.localPart, r.URL.Path)
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

			got, err := c.GetMailbox(tt.domain, tt.localPart)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMailbox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMailbox() got = %v, want %v", got, tt.want)
			}
		})
	}
}
