//go:build simulator

/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package simulator

import (
	"encoding/json"
	"fmt"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

var rewritesPattern = regexp.MustCompile("/domains/(.*)/rewrites/?(.*)?")

func handleRewrites(t *testing.T, rewrites *[]model.Rewrite, forcedStatusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := rewritesPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			t.Errorf("Expected to request to match %s, got: %s", rewritesPattern, r.URL.Path)
		}
		domain := matches[1]
		slug := matches[2]

		if forcedStatusCode > 0 {
			w.WriteHeader(forcedStatusCode)
			return
		}

		if r.Method == http.MethodPost {
			handleCreateRewrite(w, r, t, rewrites, domain)
		}
		if r.Method == http.MethodPut {
			handleUpdateRewrite(w, r, t, rewrites, domain, slug)
		}
		if r.Method == http.MethodDelete {
			handleDeleteRewrite(w, r, t, rewrites, domain, slug)
		}
		if r.Method == http.MethodGet && slug != "" {
			handleGetRewrite(w, r, t, rewrites, domain, slug)
		}
		if r.Method == http.MethodGet && slug == "" {
			handleGetRewrites(w, r, t, rewrites, domain)
		}
	}
}

func handleGetRewrites(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/rewrites", domain) {
		t.Errorf("Expected to request '/domains/%s/rewrites', got: %s", domain, r.URL.Path)
	}

	var found []model.Rewrite
	for _, rewrite := range *rewrites {
		if rewrite.DomainName == domain {
			found = append(found, rewrite)
		}
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(t, w, model.Rewrites{Rewrites: found})
}

func handleGetRewrite(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string, slug string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/rewrites/%s", domain, slug) {
		t.Errorf("Expected to request '/domains/%s/rewrites/%s', got: %s", domain, slug, r.URL.Path)
	}

	missing := true
	for _, rewrite := range *rewrites {
		if rewrite.DomainName == domain && rewrite.Name == slug {
			missing = false
			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, rewrite)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleDeleteRewrite(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string, slug string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/rewrites/%s", domain, slug) {
		t.Errorf("Expected to request '/domains/%s/rewrites/%s', got: %s", domain, slug, r.URL.Path)
	}

	missing := true
	for index, rewrite := range *rewrites {
		if rewrite.DomainName == domain && rewrite.Name == slug {
			missing = false
			c := *rewrites
			c[index] = c[len(c)-1]
			*rewrites = c[:len(c)-1]

			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, rewrite)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleUpdateRewrite(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string, slug string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/rewrites/%s", domain, slug) {
		t.Errorf("Expected to request '/domains/%s/rewrites/%s', got: %s", domain, slug, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	requestRewrite := RewriteServerModel{}
	err = json.Unmarshal(requestBody, &requestRewrite)
	if err != nil {
		t.Errorf("Could not unmarshall rewrite rule")
	}

	requestRewrite.DomainName = domain
	requestRewrite.Name = slug

	missing := true
	for index, rewrite := range *rewrites {
		if rewrite.DomainName == domain && rewrite.Name == slug {
			missing = false
			responseRewrite := model.Rewrite{
				DomainName:    domain,
				Name:          requestRewrite.Name,
				LocalPartRule: requestRewrite.LocalPartRule,
				OrderNum:      requestRewrite.OrderNum,
				Destinations:  strings.Split(requestRewrite.Destinations, ","),
			}
			c := *rewrites
			c[index] = responseRewrite
			*rewrites = c

			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, responseRewrite)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleCreateRewrite(w http.ResponseWriter, r *http.Request, t *testing.T, rewrites *[]model.Rewrite, domain string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/rewrites", domain) {
		t.Errorf("Expected to request '/domains/%s/rewrites', got: %s", domain, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	requestRewrite := RewriteServerModel{}
	err = json.Unmarshal(requestBody, &requestRewrite)
	if err != nil {
		t.Errorf("Could not unmarshall rewrite rule")
	}

	for _, existingRewrite := range *rewrites {
		if existingRewrite.DomainName == domain && existingRewrite.Name == requestRewrite.Name {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	rewrite := model.Rewrite{
		DomainName:    domain,
		Name:          requestRewrite.Name,
		LocalPartRule: requestRewrite.LocalPartRule,
		OrderNum:      requestRewrite.OrderNum,
		Destinations:  strings.Split(requestRewrite.Destinations, ","),
	}

	*rewrites = append(*rewrites, rewrite)

	w.WriteHeader(http.StatusOK)
	writeJsonResponse(t, w, rewrite)
}

type RewriteServerModel struct {
	model.Rewrite
	Destinations string `json:"destinations"`
}
