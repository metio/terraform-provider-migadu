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

func TestMigaduClient_GetAliases(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		state      []model.Alias
		want       *model.Aliases
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Destinations:     []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Aliases{
				Aliases: []model.Alias{
					{
						LocalPart:        "some",
						DomainName:       "example.com",
						Address:          "some@example.com",
						Destinations:     []string{"other@example"},
						IsInternal:       true,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: false,
					},
				},
			},
		},
		{
			name:   "multiple",
			domain: "example.com",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Destinations:     []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
				{
					LocalPart:        "other",
					DomainName:       "example.com",
					Address:          "other@example.com",
					Destinations:     []string{"different@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Aliases{
				Aliases: []model.Alias{
					{
						LocalPart:        "some",
						DomainName:       "example.com",
						Address:          "some@example.com",
						Destinations:     []string{"other@example"},
						IsInternal:       true,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: false,
					},
					{
						LocalPart:        "other",
						DomainName:       "example.com",
						Address:          "other@example.com",
						Destinations:     []string{"different@example"},
						IsInternal:       true,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: false,
					},
				},
			},
		},
		{
			name:   "filtered",
			domain: "example.com",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "other.com",
					Address:          "some@other.com",
					Destinations:     []string{"other@other"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
				{
					LocalPart:        "other",
					DomainName:       "example.com",
					Address:          "other@example.com",
					Destinations:     []string{"different@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Aliases{
				Aliases: []model.Alias{
					{
						LocalPart:        "other",
						DomainName:       "example.com",
						Address:          "other@example.com",
						Destinations:     []string{"different@example"},
						IsInternal:       true,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: false,
					},
				},
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			state: []model.Alias{
				{
					LocalPart:        "test",
					DomainName:       "xn--ho-hia.de",
					Address:          "test@xn--ho-hia.de",
					Destinations:     []string{"other@xn--ho-hia.de"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Aliases{
				Aliases: []model.Alias{
					{
						LocalPart:        "test",
						DomainName:       "xn--ho-hia.de",
						Address:          "test@xn--ho-hia.de",
						Destinations:     []string{"other@xn--ho-hia.de"},
						IsInternal:       true,
						Expirable:        false,
						ExpiresOn:        "",
						RemoveUponExpiry: false,
					},
				},
			},
		},
		{
			name:       "error-400",
			domain:     "hoß.de",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "error-404",
			domain:     "hoß.de",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "error-500",
			domain:     "hoß.de",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Aliases: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetAliases(context.Background(), tt.domain)
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
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Alias
		want       *model.Alias
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "some",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Destinations:     []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:        "some",
				DomainName:       "example.com",
				Address:          "some@example.com",
				Destinations:     []string{"other@example"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "multiple",
			domain:    "example.com",
			localPart: "some",
			state: []model.Alias{
				{
					LocalPart:        "some",
					DomainName:       "example.com",
					Address:          "some@example.com",
					Destinations:     []string{"other@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
				{
					LocalPart:        "other",
					DomainName:       "example.com",
					Address:          "other@example.com",
					Destinations:     []string{"different@example"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:        "some",
				DomainName:       "example.com",
				Address:          "some@example.com",
				Destinations:     []string{"other@example"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:        "test",
					DomainName:       "xn--ho-hia.de",
					Address:          "test@xn--ho-hia.de",
					Destinations:     []string{"other@xn--ho-hia.de"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:        "test",
				DomainName:       "xn--ho-hia.de",
				Address:          "test@xn--ho-hia.de",
				Destinations:     []string{"other@xn--ho-hia.de"},
				IsInternal:       true,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:       "error-400",
			domain:     "hoß.de",
			localPart:  "test",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:      "error-404",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:        "test",
					DomainName:       "example.com",
					Address:          "test@example.com",
					Destinations:     []string{"other@example.com"},
					IsInternal:       true,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "hoß.de",
			localPart:  "test",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Aliases: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetAlias(context.Background(), tt.domain, tt.localPart)
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

func TestMigaduClient_CreateAlias(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		state      []model.Alias
		send       *model.Alias
		want       *model.Alias
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			state:  []model.Alias{},
			send: &model.Alias{
				LocalPart: "other",
				Destinations: []string{
					"different@example.com",
				},
			},
			want: &model.Alias{
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
		},
		{
			name:   "idna",
			domain: "hoß.de",
			state:  []model.Alias{},
			send: &model.Alias{
				LocalPart: "test",
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
			want: &model.Alias{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
		},
		{
			name:   "error-duplicate",
			domain: "example.com",
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"another@example.com",
					},
				},
			},
			send: &model.Alias{
				LocalPart: "test",
				Destinations: []string{
					"another@example.com",
				},
			},
			wantErr: true,
		},
		{
			name:   "error-duplicate-idna",
			domain: "hoß.de",
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Destinations: []string{
						"another@xn--ho-hia.de",
					},
				},
			},
			send: &model.Alias{
				LocalPart: "test",
				Destinations: []string{
					"another@hoß.de",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Aliases: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.CreateAlias(context.Background(), tt.domain, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAlias() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_UpdateAlias(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Alias
		send       *model.Alias
		want       *model.Alias
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"other@example.com",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			send: &model.Alias{
				Destinations: []string{
					"different@example.com",
				},
			},
			want: &model.Alias{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Destinations: []string{
					"different@example.com",
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
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"other@example.com",
						"another@example.com",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			send: &model.Alias{
				Destinations: []string{
					"different@example.com",
				},
			},
			want: &model.Alias{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
				Destinations: []string{
					"different@example.com",
				},
				IsInternal:       false,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Destinations: []string{
						"other@xn--ho-hia.de",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			send: &model.Alias{
				Destinations: []string{
					"different@hoß.de",
				},
			},
			want: &model.Alias{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Destinations: []string{
					"different@xn--ho-hia.de",
				},
				IsInternal:       false,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "another.com",
					Address:    "test@another.com",
					Destinations: []string{
						"other@example.com",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			send: &model.Alias{
				Destinations: []string{
					"different@example.com",
				},
			},
			wantErr: true,
		},
		{
			name:      "error-500",
			domain:    "example.com",
			localPart: "test",
			state:     []model.Alias{},
			send: &model.Alias{
				Destinations: []string{
					"different@example.com",
				},
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Aliases: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.UpdateAlias(context.Background(), tt.domain, tt.localPart, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateAlias() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_DeleteAlias(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		localPart  string
		statusCode int
		state      []model.Alias
		want       *model.Alias
		wantErr    bool
	}{
		{
			name:      "single",
			domain:    "example.com",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"other@example.com",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
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
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "example.com",
					Address:    "test@example.com",
					Destinations: []string{
						"other@example.com",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
				{
					LocalPart:  "test",
					DomainName: "another.com",
					Address:    "test@another.com",
					Destinations: []string{
						"other@another.com",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:  "test",
				DomainName: "example.com",
				Address:    "test@example.com",
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
			name:      "idna",
			domain:    "hoß.de",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "xn--ho-hia.de",
					Address:    "test@xn--ho-hia.de",
					Destinations: []string{
						"other@xn--ho-hia.de",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
				},
			},
			want: &model.Alias{
				LocalPart:  "test",
				DomainName: "xn--ho-hia.de",
				Address:    "test@xn--ho-hia.de",
				Destinations: []string{
					"other@xn--ho-hia.de",
				},
				IsInternal:       false,
				Expirable:        false,
				ExpiresOn:        "",
				RemoveUponExpiry: false,
			},
		},
		{
			name:      "error-404",
			domain:    "example.com",
			localPart: "test",
			state: []model.Alias{
				{
					LocalPart:  "test",
					DomainName: "another.com",
					Address:    "test@another.com",
					Destinations: []string{
						"other@another.com",
					},
					IsInternal:       false,
					Expirable:        false,
					ExpiresOn:        "",
					RemoveUponExpiry: false,
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
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Aliases: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.DeleteAlias(context.Background(), tt.domain, tt.localPart)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteAlias() got = %v, want %v", got, tt.want)
			}
		})
	}
}
