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

func TestMigaduClient_GetIdentities(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Identity
		want       *model.Identities
		wantErr    bool
	}{
		{
			name:      "empty",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Identity{},
			want:      &model.Identities{},
		},
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:            "other",
					DomainName:           "example.com",
					Address:              "other@example.com",
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
			want: &model.Identities{
				Identities: []model.Identity{
					{
						LocalPart:            "other",
						DomainName:           "example.com",
						Address:              "other@example.com",
						Name:                 "Some Name",
						MaySend:              true,
						MayReceive:           true,
						MayAccessImap:        true,
						MayAccessPop3:        true,
						MayAccessManageSieve: true,
						FooterActive:         false,
						FooterPlainBody:      "",
						FooterHtmlBody:       "",
					},
				},
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "another",
					DomainName: "example.com",
					Address:    "another@example.com",
					Name:       "Another Name",
				},
			},
			want: &model.Identities{
				Identities: []model.Identity{
					{
						LocalPart:  "other",
						DomainName: "example.com",
						Address:    "other@example.com",
						Name:       "Some Name",
					},
					{
						LocalPart:  "another",
						DomainName: "example.com",
						Address:    "another@example.com",
						Name:       "Another Name",
					},
				},
			},
		},
		{
			name:      "filtered",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "another",
					DomainName: "different.com",
					Address:    "another@different.com",
					Name:       "Another Name",
				},
			},
			want: &model.Identities{
				Identities: []model.Identity{
					{
						LocalPart:  "other",
						DomainName: "example.com",
						Address:    "other@example.com",
						Name:       "Some Name",
					},
				},
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "xn--ho-hia.de",
					Address:    "other@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: &model.Identities{
				Identities: []model.Identity{
					{
						LocalPart:  "other",
						DomainName: "xn--ho-hia.de",
						Address:    "other@xn--ho-hia.de",
						Name:       "Some Name",
					},
				},
			},
		},
		{
			name:       "error-404",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusNotFound,
			wantErr:    true,
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
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetIdentities(context.Background(), tt.domain, tt.localPart)
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
		name       string
		domain     string
		localPart  string
		identity   string
		statusCode int
		state      []model.Identity
		want       *model.Identity
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "xn--ho-hia.de",
					Address:    "other@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "xn--ho-hia.de",
				Address:    "other@xn--ho-hia.de",
				Name:       "Some Name",
			},
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "different",
					DomainName: "example.com",
					Address:    "different@example.com",
					Name:       "Different Name",
				},
			},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			identity:   "other",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetIdentity(context.Background(), tt.domain, tt.localPart, tt.identity)
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

func TestMigaduClient_CreateIdentity(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Identity
		send       *model.Identity
		want       *model.Identity
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Identity{},
			send: &model.Identity{
				LocalPart: "other",
				Name:      "Some Name",
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "different",
					DomainName: "example.com",
					Address:    "different@example.com",
					Name:       "Different Name",
				},
			},
			send: &model.Identity{
				LocalPart: "other",
				Name:      "Some Name",
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state:     []model.Identity{},
			send: &model.Identity{
				LocalPart: "other",
				Name:      "Some Name",
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "xn--ho-hia.de",
				Address:    "other@xn--ho-hia.de",
				Name:       "Some Name",
			},
		},
		{
			name:      "error-400",
			domain:    "example.com",
			localPart: "test",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
			},
			send: &model.Identity{
				LocalPart: "other",
				Name:      "Some Name",
			},
			wantErr: true,
		},
		{
			name:       "error-404",
			domain:     "example.com",
			localPart:  "test",
			statusCode: http.StatusNotFound,
			wantErr:    true,
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
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.CreateIdentity(context.Background(), tt.domain, tt.localPart, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateIdentity() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_UpdateIdentity(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		identity   string
		statusCode int
		state      []model.Identity
		send       *model.Identity
		want       *model.Identity
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
			},
			send: &model.Identity{
				Name: "Different Name",
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Different Name",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "another",
					DomainName: "example.com",
					Address:    "another@example.com",
					Name:       "Another Name",
				},
			},
			send: &model.Identity{
				Name: "Different Name",
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Different Name",
			},
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			identity:  "not-found",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
			},
			send: &model.Identity{
				Name: "Different Name",
			},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			identity:   "other",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.UpdateIdentity(context.Background(), tt.domain, tt.localPart, tt.identity, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateIdentity() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_DeleteIdentity(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		identity   string
		statusCode int
		state      []model.Identity
		want       *model.Identity
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "test",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
				{
					LocalPart:  "another",
					DomainName: "example.com",
					Address:    "different@example.com",
					Name:       "Some Name",
				},
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "example.com",
				Address:    "other@example.com",
				Name:       "Some Name",
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			identity:  "other",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "xn--ho-hia.de",
					Address:    "other@xn--ho-hia.de",
					Name:       "Some Name",
				},
			},
			want: &model.Identity{
				LocalPart:  "other",
				DomainName: "xn--ho-hia.de",
				Address:    "other@xn--ho-hia.de",
				Name:       "Some Name",
			},
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			identity:  "not-found",
			state: []model.Identity{
				{
					LocalPart:  "other",
					DomainName: "example.com",
					Address:    "other@example.com",
					Name:       "Some Name",
				},
			},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			localPart:  "test",
			identity:   "server-error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Identities: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.DeleteIdentity(context.Background(), tt.domain, tt.localPart, tt.identity)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteIdentity() got = %v, want %v", got, tt.want)
			}
		})
	}
}
