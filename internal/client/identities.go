/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/idna"
	"net/http"
)

type Identities struct {
	Identities []Identity `json:"identities"`
}

type Identity struct {
	LocalPart            string `json:"local_part"`
	DomainName           string `json:"domain_name"`
	Address              string `json:"address"`
	Name                 string `json:"name"`
	MaySend              bool   `json:"may_send"`
	MayReceive           bool   `json:"may_receive"`
	MayAccessImap        bool   `json:"may_access_imap"`
	MayAccessPop3        bool   `json:"may_access_pop3"`
	MayAccessManageSieve bool   `json:"may_access_managesieve"`
	Password             string `json:"password"`
	FooterActive         bool   `json:"footer_active"`
	FooterPlainBody      string `json:"footer_plain_body"`
	FooterHtmlBody       string `json:"footer_html_body"`
}

// GetIdentities - Returns identities for a single mailbox
func (c *MigaduClient) GetIdentities(ctx context.Context, domain string, localPart string) (*Identities, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s/identities", c.Endpoint, ascii, localPart)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	response := Identities{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	return &response, nil
}

// GetIdentity - Returns a specific identity
func (c *MigaduClient) GetIdentity(ctx context.Context, domain string, localPart string, id string) (*Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s/identities/%s", c.Endpoint, ascii, localPart, id)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	response := Identity{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	return &response, nil
}

// CreateIdentity - Creates a new identity
func (c *MigaduClient) CreateIdentity(ctx context.Context, domain string, localPart string, identity *Identity) (*Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s/identities", c.Endpoint, ascii, localPart)

	requestBody, err := json.Marshal(identity)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	response := Identity{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	return &response, nil
}

// UpdateIdentity - Updates an existing identity
func (c *MigaduClient) UpdateIdentity(ctx context.Context, domain string, localPart string, id string, identity *Identity) (*Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s/identities/%s", c.Endpoint, ascii, localPart, id)

	requestBody, err := json.Marshal(identity)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	response := Identity{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	return &response, nil
}

// DeleteIdentity - Deletes an existing identity
func (c *MigaduClient) DeleteIdentity(ctx context.Context, domain string, localPart string, id string) (*Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteIdentity: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s/identities/%s", c.Endpoint, ascii, localPart, id)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("DeleteIdentity: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("DeleteIdentity: %w", err)
	}

	response := Identity{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("DeleteIdentity: %w", err)
	}

	return &response, nil
}
