/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package idn

import (
	"reflect"
	"testing"
)

func TestConvertEmailsToASCII(t *testing.T) {
	tests := []struct {
		name    string
		emails  []string
		want    []string
		wantErr bool
	}{
		{
			name: "single",
			emails: []string{
				"test@example.com",
			},
			want: []string{
				"test@example.com",
			},
			wantErr: false,
		},
		{
			name: "multiple",
			emails: []string{
				"test@example.com",
				"test@another.com",
			},
			want: []string{
				"test@example.com",
				"test@another.com",
			},
			wantErr: false,
		},
		{
			name: "idna",
			emails: []string{
				"test@hoß.de",
			},
			want: []string{
				"test@xn--ho-hia.de",
			},
			wantErr: false,
		},
		{
			name: "mixed",
			emails: []string{
				"test@hoß.de",
				"test@example.com",
			},
			want: []string{
				"test@xn--ho-hia.de",
				"test@example.com",
			},
			wantErr: false,
		},
		{
			name: "punycode",
			emails: []string{
				"test@xn--ho-hia.de",
			},
			want: []string{
				"test@xn--ho-hia.de",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertEmailsToASCII(tt.emails)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertEmailsToASCII() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertEmailsToASCII() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertEmailsToUnicode(t *testing.T) {
	tests := []struct {
		name    string
		emails  []string
		want    []string
		wantErr bool
	}{
		{
			name: "single",
			emails: []string{
				"test@example.com",
			},
			want: []string{
				"test@example.com",
			},
			wantErr: false,
		},
		{
			name: "multiple",
			emails: []string{
				"test@example.com",
				"test@another.com",
			},
			want: []string{
				"test@example.com",
				"test@another.com",
			},
			wantErr: false,
		},
		{
			name: "idna",
			emails: []string{
				"test@xn--ho-hia.de",
			},
			want: []string{
				"test@hoß.de",
			},
			wantErr: false,
		},
		{
			name: "mixed",
			emails: []string{
				"test@xn--ho-hia.de",
				"test@example.com",
			},
			want: []string{
				"test@hoß.de",
				"test@example.com",
			},
			wantErr: false,
		},
		{
			name: "unicode",
			emails: []string{
				"test@hoß.de",
			},
			want: []string{
				"test@hoß.de",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertEmailsToUnicode(tt.emails)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertEmailsToUnicode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertEmailsToUnicode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
