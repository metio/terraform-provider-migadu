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

var mailboxesPattern = regexp.MustCompile("/domains/(.*)/mailboxes/?(.*)?")

func handleMailboxes(t *testing.T, mailboxes *[]model.Mailbox) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := identitiesUrlPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			t.Errorf("Expected to request to match %s, got: %s", mailboxesPattern, r.URL.Path)
		}
		domain := matches[1]
		localPart := matches[2]

		if r.Method == http.MethodPost {
			handleCreateMailbox(w, r, t, mailboxes, domain)
		}
		if r.Method == http.MethodPut {
			handleUpdateMailbox(w, r, t, mailboxes, domain, localPart)
		}
		if r.Method == http.MethodDelete {
			handleDeleteMailbox(w, r, t, mailboxes, domain, localPart)
		}
		if r.Method == http.MethodGet {
			handleGetMailbox(w, r, t, mailboxes, domain, localPart)
		}
	}
}

func handleGetMailbox(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string, localPart string) {

}

func handleDeleteMailbox(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string, localPart string) {

}

func handleUpdateMailbox(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string, localPart string) {

}

func handleCreateMailbox(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string) {

}
