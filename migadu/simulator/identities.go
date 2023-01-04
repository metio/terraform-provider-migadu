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

var identitiesUrlPattern = regexp.MustCompile("/domains/(.*)/mailboxes/(.*)/identities/?(.*)?")

func handleIdentities(t *testing.T, identities *[]model.Identity, forcedStatusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := identitiesUrlPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			t.Errorf("Expected to request to match %s, got: %s", identitiesUrlPattern, r.URL.Path)
		}
		domain := matches[1]
		localPart := matches[2]
		id := matches[3]

		if forcedStatusCode > 0 {
			w.WriteHeader(forcedStatusCode)
			return
		}

		if r.Method == http.MethodPost {
			handleCreateIdentity(w, r, t, identities, domain, localPart)
		}
		if r.Method == http.MethodPut {
			handleUpdateIdentity(w, r, t, identities, domain, localPart, id)
		}
		if r.Method == http.MethodDelete {
			handleDeleteIdentity(w, r, t, identities, domain, localPart, id)
		}
		if r.Method == http.MethodGet && id != "" {
			handleGetIdentity(w, r, t, identities, domain, localPart, id)
		}
		if r.Method == http.MethodGet && id == "" {
			handleGetIdentities(w, r, t, identities, domain, localPart)
		}
	}
}

func handleGetIdentities(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s/identities", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/mailboxes/%s/identities', got: %s", domain, localPart, r.URL.Path)
	}

	var found []model.Identity
	for _, identity := range *identities {
		if identity.DomainName == domain {
			found = append(found, identityResponse(identity))
		}
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(t, w, model.Identities{Identities: found})
}

func handleGetIdentity(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string, localPart string, id string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s/identities/%s", domain, localPart, id) {
		t.Errorf("Expected to request '/domains/%s/mailboxes/%s/identities/%s', got: %s", domain, localPart, id, r.URL.Path)
	}

	missing := true
	for _, identity := range *identities {
		if identity.DomainName == domain && identity.LocalPart == id {
			missing = false
			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, identityResponse(identity))
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleDeleteIdentity(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string, localPart string, id string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s/identities/%s", domain, localPart, id) {
		t.Errorf("Expected to request '/domains/%s/mailboxes/%s/identities/%s', got: %s", domain, localPart, id, r.URL.Path)
	}

	missing := true
	for index, identity := range *identities {
		if identity.DomainName == domain && identity.LocalPart == id {
			missing = false
			c := *identities
			c[index] = c[len(c)-1]
			*identities = c[:len(c)-1]

			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, identityResponse(identity))
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleUpdateIdentity(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string, localPart string, id string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s/identities/%s", domain, localPart, id) {
		t.Errorf("Expected to request '/domains/%s/mailboxes/%s/identities/%s', got: %s", domain, localPart, id, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	requestIdentity := model.Identity{}
	err = json.Unmarshal(requestBody, &requestIdentity)
	if err != nil {
		t.Errorf("Could not unmarshall identity")
	}

	requestIdentity.DomainName = domain
	requestIdentity.LocalPart = id
	requestIdentity.Address = fmt.Sprintf("%s@%s", requestIdentity.LocalPart, domain)

	missing := true
	for index, identity := range *identities {
		if identity.DomainName == domain && identity.LocalPart == id {
			missing = false
			c := *identities
			c[index] = requestIdentity
			*identities = c

			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, identityResponse(requestIdentity))
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleCreateIdentity(w http.ResponseWriter, r *http.Request, t *testing.T, identities *[]model.Identity, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/mailboxes/%s/identities", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/mailboxes/%s/identities', got: %s", domain, localPart, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	identity := model.Identity{}
	err = json.Unmarshal(requestBody, &identity)
	if err != nil {
		t.Errorf("Could not unmarshall identity")
	}

	identity.DomainName = domain
	identity.Address = fmt.Sprintf("%s@%s", identity.LocalPart, domain)

	for _, existingIdentity := range *identities {
		if existingIdentity.DomainName == identity.DomainName && existingIdentity.LocalPart == identity.LocalPart {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	*identities = append(*identities, identity)

	w.WriteHeader(http.StatusOK)
	writeJsonResponse(t, w, identityResponse(identity))
}

func identityResponse(identity model.Identity) model.Identity {
	identity.Password = ""
	return identity
}
