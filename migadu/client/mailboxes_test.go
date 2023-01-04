//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package client_test

import (
	"context"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"github.com/metio/terraform-provider-migadu/migadu/simulator"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMigaduClient_GetMailboxes(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		state      []model.Mailbox
		want       *model.Mailboxes
		wantErr    bool
	}{
		{
			name:   "empty",
			domain: "example.com",
			want:   &model.Mailboxes{},
		},
		{
			name:   "single",
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "test",
				},
			},
			want: &model.Mailboxes{
				Mailboxes: []model.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "example.com",
						Address:    "test@example.com",
						Name:       "test",
					},
				},
			},
		},
		{
			name:   "multiple",
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "test",
				},
				{
					LocalPart:  "another",
					DomainName: "example.com",
					Address:    "another@example.com",
				},
			},
			want: &model.Mailboxes{
				Mailboxes: []model.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "example.com",
						Address:    "test@example.com",
						Name:       "test",
					},
					{
						LocalPart:  "another",
						DomainName: "example.com",
						Address:    "another@example.com",
					},
				},
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "test",
				},
			},
			want: &model.Mailboxes{
				Mailboxes: []model.Mailbox{
					{
						LocalPart:  "test",
						DomainName: "xn--ho-hia.de",
						Address:    "test@xn--ho-hia.de",
						Name:       "test",
					},
				},
			},
		},
		{
			name:       "error-404",
			domain:     "example.com",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetMailboxes(context.Background(), tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMailboxes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMailboxes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_GetMailbox(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Mailbox
		want       *model.Mailbox
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "different.com",
					Address:    "test@different.com",
					Name:       "Different Name",
				},
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Name:       "Some Name",
			},
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Mailbox{},
			wantErr:   true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetMailbox(context.Background(), tt.domain, tt.localPart)
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

func TestMigaduClient_CreateMailbox(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		state      []model.Mailbox
		send       *model.Mailbox
		want       *model.Mailbox
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			state:  []model.Mailbox{},
			send: &model.Mailbox{
				Name: "Some Name",
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:   "multple",
			domain: "example.com",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "different.com",
					Address:    "test@different.com",
					Name:       "Some Name",
				},
			},
			send: &model.Mailbox{
				Name: "Some Name",
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			state:  []model.Mailbox{},
			send: &model.Mailbox{
				Name: "Some Name",
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Name:       "Some Name",
			},
		},
		{
			name:       "error-404",
			domain:     "example.com",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.CreateMailbox(context.Background(), tt.domain, tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMailbox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMailbox() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_UpdateMailbox(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Mailbox
		send       *model.Mailbox
		want       *model.Mailbox
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
				},
			},
			send: &model.Mailbox{
				Name: "Different Name",
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Different Name",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "test",
					DomainName: "another.com",
					Address:    "test@another.com",
					Name:       "Some Name",
				},
			},
			send: &model.Mailbox{
				Name: "Different Name",
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Different Name",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			send: &model.Mailbox{
				Name: "Different Name",
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Name:       "Different Name",
			},
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
			},
			send: &model.Mailbox{
				Name: "Different Name",
			},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.UpdateMailbox(context.Background(), tt.domain, tt.localPart, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMailbox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateMailbox() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_DeleteMailbox(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Mailbox
		want       *model.Mailbox
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "different",
					DomainName: "example.com",
					Address:    "different@example.com",
					Name:       "Different Name",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: &model.Mailbox{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Name:       "Some Name",
			},
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			state: []model.Mailbox{
				{
					LocalPart:  "test",
					DomainName: "different.com",
					Address:    "test@different.com",
					Name:       "Some Name",
				},
			},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Mailboxes: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.DeleteMailbox(context.Background(), tt.domain, tt.localPart)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMailbox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteMailbox() got = %v, want %v", got, tt.want)
			}
		})
	}
}
