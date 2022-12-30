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

type Aliases struct {
	Aliases []Alias `json:"address_aliases"`
}

type Alias struct {
	LocalPart        string   `json:"local_part"`
	DomainName       string   `json:"domain_name"`
	Address          string   `json:"address"`
	Destinations     []string `json:"destinations"`
	IsInternal       bool     `json:"is_internal"`
	Expirable        bool     `json:"expireable"`
	ExpiresOn        string   `json:"expires_on"`
	RemoveUponExpiry bool     `json:"remove_upon_expiry"`
}

// GetAliases - Returns aliases for a single domain
func (c *MigaduClient) GetAliases(ctx context.Context, domain string) (*Aliases, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/aliases", c.Endpoint, ascii)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}

	response := Aliases{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}

	return &response, nil
}

// GetAlias - Returns specific alias
func (c *MigaduClient) GetAlias(ctx context.Context, domain string, localPart string) (*Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/aliases/%s", c.Endpoint, ascii, localPart)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	response := Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	return &response, nil
}

// CreateAlias - Creates a new alias
func (c *MigaduClient) CreateAlias(ctx context.Context, domain string, alias *Alias) (*Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/aliases", c.Endpoint, ascii)

	requestBody, err := json.Marshal(alias)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	response := Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	return &response, nil
}

// UpdateAlias - Updates an existing alias
func (c *MigaduClient) UpdateAlias(ctx context.Context, domain string, localPart string, alias *Alias) (*Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/aliases/%s", c.Endpoint, ascii, localPart)

	requestBody, err := json.Marshal(alias)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	response := Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	return &response, nil
}

// DeleteAlias - Deletes an existing alias
func (c *MigaduClient) DeleteAlias(ctx context.Context, domain string, localPart string) (*Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/aliases/%s", c.Endpoint, ascii, localPart)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	response := Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	return &response, nil
}
