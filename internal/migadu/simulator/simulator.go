/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package simulator

import (
	"github.com/metio/terraform-provider-migadu/internal/migadu/model"
	"net/http"
	"testing"
)

type State struct {
	Mailboxes  []model.Mailbox
	Aliases    []model.Alias
	Identities []model.Identity
	Rewrites   []model.Rewrite
}

func MigaduAPI(t *testing.T, state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if aliasesUrlPattern.MatchString(r.URL.Path) {
			handleAliases(t, &state.Aliases).ServeHTTP(w, r)
		} else if identitiesUrlPattern.MatchString(r.URL.Path) {
			handleIdentities(t, &state.Identities).ServeHTTP(w, r)
		} else if rewritesPattern.MatchString(r.URL.Path) {
			handleRewrites(t, &state.Rewrites).ServeHTTP(w, r)
		} else if mailboxesPattern.MatchString(r.URL.Path) {
			handleMailboxes(t, &state.Mailboxes).ServeHTTP(w, r)
		} else {
			t.Errorf("No Handler for URL: %s", r.URL.Path)
		}
	}
}
