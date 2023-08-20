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
		state      []model.RewriteRule
		want       *model.RewriteRules
		wantErr    bool
	}{
		{
			name:   "empty",
			domain: "example.com",
			want:   &model.RewriteRules{},
		},
		{
			name:   "single",
			domain: "example.com",
			state: []model.RewriteRule{
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
			want: &model.RewriteRules{
				RewriteRules: []model.RewriteRule{
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
			state: []model.RewriteRule{
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
			want: &model.RewriteRules{
				RewriteRules: []model.RewriteRule{
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
			state: []model.RewriteRule{
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
			want: &model.RewriteRules{
				RewriteRules: []model.RewriteRule{
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

			got, err := c.GetRewriteRules(context.Background(), tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRewriteRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRewriteRules() got = %v, want %v", got, tt.want)
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
		state      []model.RewriteRule
		want       *model.RewriteRule
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			slug:   "slug",
			state: []model.RewriteRule{
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
			want: &model.RewriteRule{
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
			state: []model.RewriteRule{
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
			want: &model.RewriteRule{
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
			state: []model.RewriteRule{
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
			want: &model.RewriteRule{
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
			state:   []model.RewriteRule{},
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

			got, err := c.GetRewriteRule(context.Background(), tt.domain, tt.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRewriteRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRewriteRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigaduClient_CreateRewrite(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		statusCode int
		state      []model.RewriteRule
		send       *model.RewriteRule
		want       *model.RewriteRule
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			state: []model.RewriteRule{
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
			send: &model.RewriteRule{
				Name:          "sec",
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"security@example.com",
				},
			},
			want: &model.RewriteRule{
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
			state: []model.RewriteRule{
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
			send: &model.RewriteRule{
				Name:          "slug",
				LocalPartRule: "rule-*",
				OrderNum:      1,
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
			want: &model.RewriteRule{
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
			state: []model.RewriteRule{
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
			send: &model.RewriteRule{
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

			got, err := c.CreateRewriteRule(context.Background(), tt.domain, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRewriteRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRewriteRule() got = %v, want %v", got, tt.want)
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
		state      []model.RewriteRule
		send       *model.RewriteRule
		want       *model.RewriteRule
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			slug:   "slug",
			state: []model.RewriteRule{
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
			send: &model.RewriteRule{
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"another@example.com",
				},
			},
			want: &model.RewriteRule{
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
			state: []model.RewriteRule{
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
			send: &model.RewriteRule{
				LocalPartRule: "sec-*",
				OrderNum:      0,
				Destinations: []string{
					"another@xn--ho-hia.de",
				},
			},
			want: &model.RewriteRule{
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
			state: []model.RewriteRule{
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
			send: &model.RewriteRule{
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

			got, err := c.UpdateRewriteRule(context.Background(), tt.domain, tt.slug, tt.send)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRewriteRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateRewriteRule() got = %v, want %v", got, tt.want)
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
		state      []model.RewriteRule
		want       *model.RewriteRule
		wantErr    bool
	}{
		{
			name:   "single",
			domain: "example.com",
			slug:   "slug",
			state: []model.RewriteRule{
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
			want: &model.RewriteRule{
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
			state: []model.RewriteRule{
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
			want: &model.RewriteRule{
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
			state: []model.RewriteRule{
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
			want: &model.RewriteRule{
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
			state:   []model.RewriteRule{},
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

			got, err := c.DeleteRewriteRule(context.Background(), tt.domain, tt.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRewriteRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRewriteRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}
