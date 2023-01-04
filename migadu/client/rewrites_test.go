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

func TestMigaduClient_GetRewrites(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		state      []model.Rewrite
		want       *model.Rewrites
		wantErr    bool
	}{
		{
			name:   "empty",
			domain: "example.com",
			want:   &model.Rewrites{},
		},
		{
			name:   "single",
			domain: "example.com",
			state: []model.Rewrite{
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
			want: &model.Rewrites{
				Rewrites: []model.Rewrite{
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
		},
		{
			name:   "multiple",
			domain: "example.com",
			state: []model.Rewrite{
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
					DomainName:    "different.com",
					Name:          "test",
					LocalPartRule: "rule-*",
					OrderNum:      1,
					Destinations: []string{
						"another@different.com",
					},
				},
			},
			want: &model.Rewrites{
				Rewrites: []model.Rewrite{
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
		},
		{
			name:   "idna",
			domain: "hoß.de",
			state: []model.Rewrite{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "test",
					LocalPartRule: "rule-*",
					OrderNum:      1,
					Destinations: []string{
						"another@xn--ho-hia.de",
					},
				},
			},
			want: &model.Rewrites{
				Rewrites: []model.Rewrite{
					{
						DomainName:    "xn--ho-hia.de",
						Name:          "test",
						LocalPartRule: "rule-*",
						OrderNum:      1,
						Destinations: []string{
							"another@xn--ho-hia.de",
						},
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
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetRewrites(context.Background(), tt.domain)
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
		name       string
		domain     string
		slug       string
		statusCode int
		state      []model.Rewrite
		want       *model.Rewrite
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			slug:   "slug",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "slug",
					LocalPartRule: "rule-*",
					OrderNum:      1,
					Destinations: []string{
						"another@example.com",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "slug",
				LocalPartRule: "rule-*",
				OrderNum:      1,
				Destinations: []string{
					"another@example.com",
				},
			},
		},
		{
			name:   "multiple",
			domain: "example.com",
			slug:   "slug",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "slug",
					LocalPartRule: "rule-*",
					OrderNum:      1,
					Destinations: []string{
						"another@example.com",
					},
				},
				{
					DomainName:    "different.com",
					Name:          "slug",
					LocalPartRule: "rule-*",
					OrderNum:      1,
					Destinations: []string{
						"another@different.com",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "slug",
				LocalPartRule: "rule-*",
				OrderNum:      1,
				Destinations: []string{
					"another@example.com",
				},
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			slug:   "slug",
			state: []model.Rewrite{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "slug",
					LocalPartRule: "rule-*",
					OrderNum:      1,
					Destinations: []string{
						"another@xn--ho-hia.de",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "xn--ho-hia.de",
				Name:          "slug",
				LocalPartRule: "rule-*",
				OrderNum:      1,
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
		},
		{
			name:    "error-404",
			domain:  "example.com",
			slug:    "slug",
			state:   []model.Rewrite{},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			slug:       "slug",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.GetRewrite(context.Background(), tt.domain, tt.slug)
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

func TestMigaduClient_CreateRewrite(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		state      []model.Rewrite
		send       *model.Rewrite
		want       *model.Rewrite
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			state: []model.Rewrite{
				{
					DomainName:    "different.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@different.com",
					},
				},
			},
			send: &model.Rewrite{
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			send: &model.Rewrite{
				Name:          "slug",
				LocalPartRule: "rule-*",
				OrderNum:      1,
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
			want: &model.Rewrite{
				DomainName:    "xn--ho-hia.de",
				Name:          "slug",
				LocalPartRule: "rule-*",
				OrderNum:      1,
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
		},
		{
			name:   "error-404",
			domain: "example.com",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			send: &model.Rewrite{
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
			wantErr: true,
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
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.CreateRewrite(context.Background(), tt.domain, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRewrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRewrite() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_UpdateRewrite(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		slug       string
		statusCode int
		state      []model.Rewrite
		send       *model.Rewrite
		want       *model.Rewrite
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			slug:   "slug",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "slug",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			send: &model.Rewrite{
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"another@example.com",
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "slug",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"another@example.com",
				},
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			slug:   "sec",
			state: []model.Rewrite{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "sec",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@xn--ho-hia.de",
					},
				},
			},
			send: &model.Rewrite{
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
			want: &model.Rewrite{
				DomainName:    "xn--ho-hia.de",
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
		},
		{
			name:   "error-404",
			domain: "example.com",
			slug:   "slug",
			state: []model.Rewrite{
				{
					DomainName:    "different.com",
					Name:          "slug",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@different.com",
					},
				},
			},
			send: &model.Rewrite{
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			slug:       "slug",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.UpdateRewrite(context.Background(), tt.domain, tt.slug, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRewrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateRewrite() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_DeleteRewrite(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		slug       string
		statusCode int
		state      []model.Rewrite
		want       *model.Rewrite
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			slug:   "slug",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "slug",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@example.com",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "slug",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
		},
		{
			name:   "multiple",
			domain: "example.com",
			slug:   "slug",
			state: []model.Rewrite{
				{
					DomainName:    "example.com",
					Name:          "slug",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@example.com",
					},
				},
				{
					DomainName:    "different.com",
					Name:          "slug",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@different.com",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "example.com",
				Name:          "slug",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
		},
		{
			name:   "idna",
			domain: "hoß.de",
			slug:   "slug",
			state: []model.Rewrite{
				{
					DomainName:    "xn--ho-hia.de",
					Name:          "slug",
					LocalPartRule: "sec-*",
					OrderNum:      0,
					Destinations: []string{
						"security@xn--ho-hia.de",
					},
				},
			},
			want: &model.Rewrite{
				DomainName:    "xn--ho-hia.de",
				Name:          "slug",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@xn--ho-hia.de",
				},
			},
		},
		{
			name:    "error-404",
			domain:  "example.com",
			slug:    "slug",
			state:   []model.Rewrite{},
			wantErr: true,
		},
		{
			name:       "error-500",
			domain:     "example.com",
			slug:       "slug",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(simulator.MigaduAPI(t, &simulator.State{Rewrites: tt.state, StatusCode: tt.statusCode}))
			defer server.Close()

			c := newTestClient(server.URL)

			got, err := c.DeleteRewrite(context.Background(), tt.domain, tt.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRewrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRewrite() got = %v, want %v", got, tt.want)
			}
		})
	}
}
