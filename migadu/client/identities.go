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
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"golang.org/x/net/idna"
	"net/http"
)

// GetIdentities returns identities for a single mailbox
func (c *MigaduClient) GetIdentities(ctx context.Context, domain string, localPart string) (*model.Identities, error) {
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

	response := model.Identities{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}

	return &response, nil
}

// GetIdentity returns a specific identity
func (c *MigaduClient) GetIdentity(ctx context.Context, domain string, localPart string, id string) (*model.Identity, error) {
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

	response := model.Identity{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}

	return &response, nil
}

// CreateIdentity creates a new identity
func (c *MigaduClient) CreateIdentity(ctx context.Context, domain string, localPart string, identity *model.Identity) (*model.Identity, error) {
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

	response := model.Identity{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}

	return &response, nil
}

// UpdateIdentity updates an existing identity
func (c *MigaduClient) UpdateIdentity(ctx context.Context, domain string, localPart string, id string, identity *model.Identity) (*model.Identity, error) {
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

	response := model.Identity{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}

	return &response, nil
}

// DeleteIdentity deletes an existing identity
func (c *MigaduClient) DeleteIdentity(ctx context.Context, domain string, localPart string, id string) (*model.Identity, error) {
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

	response := model.Identity{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("DeleteIdentity: %w", err)
	}

	return &response, nil
}
