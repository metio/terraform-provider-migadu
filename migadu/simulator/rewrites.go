//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package simulator

import (
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"net/http"
	"regexp"
	"testing"
)

var rewritesPattern = regexp.MustCompile("/domains/(.*)/rewrites/?(.*)?")

func handleRewrites(t *testing.T, rewrites *[]model.Rewrite) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := rewritesPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			t.Errorf("Expected to request to match %s, got: %s", rewritesPattern, r.URL.Path)
		}
		domain := matches[1]
		slug := matches[2]

		if r.Method == http.MethodPost {
			handleCreateRewrite(w, r, t, rewrites, domain)
		}
		if r.Method == http.MethodPut {
			handleUpdateRewrite(w, r, t, rewrites, domain, slug)
		}
		if r.Method == http.MethodDelete {
			handleDeleteRewrite(w, r, t, rewrites, domain, slug)
		}
		if r.Method == http.MethodGet {
			handleGetRewrite(w, r, t, rewrites, domain, slug)
		}
	}
}

func handleGetRewrite(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string, slug string) {

}

func handleDeleteRewrite(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string, slug string) {

}

func handleUpdateRewrite(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string, slug string) {

}

func handleCreateRewrite(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string) {

}
