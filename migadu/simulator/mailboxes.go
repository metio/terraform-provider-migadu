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
	"testing"
)

var mailboxesPattern = regexp.MustCompile("/domains/(.*)/mailboxes/?(.*)?")

func handleMailboxes(t *testing.T, mailboxes *[]model.Mailbox, forcedStatusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := mailboxesPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			t.Errorf("Expected to request to match %s, got: %s", mailboxesPattern, r.URL.Path)
		}
		domain := matches[1]
		localPart := matches[2]

		if forcedStatusCode > 0 {
			w.WriteHeader(forcedStatusCode)
			return
		}

		if r.Method == http.MethodPost {
			handleCreateMailbox(w, r, t, mailboxes, domain)
		}
		if r.Method == http.MethodPut {
			handleUpdateMailbox(w, r, t, mailboxes, domain, localPart)
		}
		if r.Method == http.MethodDelete {
			handleDeleteMailbox(w, r, t, mailboxes, domain, localPart)
		}
		if r.Method == http.MethodGet && localPart != "" {
			handleGetMailbox(w, r, t, mailboxes, domain, localPart)
		}
		if r.Method == http.MethodGet && localPart == "" {
			handleGetMailboxes(w, r, t, mailboxes, domain)
		}
	}
}

func handleGetMailboxes(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes", domain) {
		t.Errorf("Expected to request '/domains/%s/mailboxes', got: %s", domain, r.URL.Path)
	}

	var found []model.Mailbox
	for _, mailbox := range *mailboxes {
		if mailbox.DomainName == domain {
			found = append(found, mailbox)
		}
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(t, w, model.Mailboxes{Mailboxes: found})
}

func handleGetMailbox(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/mailboxes/%s', got: %s", domain, localPart, r.URL.Path)
	}

	missing := true
	for _, mailbox := range *mailboxes {
		if mailbox.DomainName == domain && mailbox.LocalPart == localPart {
			missing = false
			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, mailbox)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleDeleteMailbox(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/mailboxes/%s', got: %s", domain, localPart, r.URL.Path)
	}

	missing := true
	for index, mailbox := range *mailboxes {
		if mailbox.DomainName == domain && mailbox.LocalPart == localPart {
			missing = false
			c := *mailboxes
			c[index] = c[len(c)-1]
			*mailboxes = c[:len(c)-1]

			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, mailbox)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleUpdateMailbox(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/mailboxes/%s', got: %s", domain, localPart, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	requestMailbox := model.Mailbox{}
	err = json.Unmarshal(requestBody, &requestMailbox)
	if err != nil {
		t.Errorf("Could not unmarshall mailbox")
	}

	if requestMailbox.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestMailbox.DomainName = domain
	requestMailbox.LocalPart = localPart
	requestMailbox.Address = fmt.Sprintf("%s@%s", requestMailbox.LocalPart, domain)
	requestMailbox.Password = ""

	missing := true
	for index, mailbox := range *mailboxes {
		if mailbox.DomainName == domain && mailbox.LocalPart == localPart {
			missing = false
			c := *mailboxes
			c[index] = requestMailbox
			*mailboxes = c

			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, requestMailbox)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleCreateMailbox(w http.ResponseWriter, r *http.Request, t *testing.T, mailboxes *[]model.Mailbox, domain string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes", domain) {
		t.Errorf("Expected to request '/domains/%s/mailboxes', got: %s", domain, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	mailbox := model.Mailbox{}
	err = json.Unmarshal(requestBody, &mailbox)
	if err != nil {
		t.Errorf("Could not unmarshall mailbox")
	}

	if mailbox.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if mailbox.PasswordMethod == "invitation" && mailbox.PasswordRecoveryEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mailbox.DomainName = domain
	mailbox.Address = fmt.Sprintf("%s@%s", mailbox.LocalPart, domain)
	mailbox.Password = ""

	for _, existingMailbox := range *mailboxes {
		if existingMailbox.DomainName == mailbox.DomainName && existingMailbox.LocalPart == mailbox.LocalPart {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	*mailboxes = append(*mailboxes, mailbox)

	w.WriteHeader(http.StatusOK)
	writeJsonResponse(t, w, mailbox)
}
