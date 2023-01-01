/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package simulator

import (
	"github.com/metio/terraform-provider-migadu/internal/migadu/model"
	"net/http"
	"regexp"
	"testing"
)

var identitiesUrlPattern = regexp.MustCompile("/domains/(.*)/mailboxes/(.*)/identities/?(.*)?")

func handleIdentities(t *testing.T, identities *[]model.Identity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := identitiesUrlPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			t.Errorf("Expected to request to match %s, got: %s", identitiesUrlPattern, r.URL.Path)
		}
		domain := matches[1]
		localPart := matches[2]
		identity := matches[3]

		if r.Method == http.MethodPost {
			handleCreateIdentity(w, r, t, identities, domain)
		}
		if r.Method == http.MethodPut {
			handleUpdateIdentity(w, r, t, identities, domain, localPart, identity)
		}
		if r.Method == http.MethodDelete {
			handleDeleteIdentity(w, r, t, identities, domain, localPart, identity)
		}
		if r.Method == http.MethodGet {
			handleGetIdentity(w, r, t, identities, domain, localPart, identity)
		}
	}
}

func handleGetIdentity(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string, localPart string, identity string) {

}

func handleDeleteIdentity(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string, localPart string, identity string) {

}

func handleUpdateIdentity(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string, localPart string, identity string) {

}

func handleCreateIdentity(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string) {

}
