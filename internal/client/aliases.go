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
func (c *MigaduClient) GetAliases(domain string) (*Aliases, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/domains/%s/aliases", c.Endpoint, ascii)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	aliases := Aliases{}
	err = json.Unmarshal(body, &aliases)
	if err != nil {
		return nil, err
	}

	return &aliases, nil
}

// GetAlias - Returns specific alias
func (c *MigaduClient) GetAlias(domain string, localPart string) (*Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/domains/%s/aliases/%s", c.Endpoint, ascii, localPart)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	alias := Alias{}
	err = json.Unmarshal(body, &alias)
	if err != nil {
		return nil, err
	}

	return &alias, nil
}

// CreateAlias - Creates a new alias
func (c *MigaduClient) CreateAlias(ctx context.Context, domain string, alias *Alias) (*Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/domains/%s/aliases", c.Endpoint, ascii)

	requestBody, err := json.Marshal(alias)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	response := Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateAlias - Updates an existing alias
func (c *MigaduClient) UpdateAlias(ctx context.Context, domain string, localPart string, alias *Alias) (*Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/domains/%s/aliases/%s", c.Endpoint, ascii, localPart)

	requestBody, err := json.Marshal(alias)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	response := Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteAlias - Deletes an existing alias
func (c *MigaduClient) DeleteAlias(ctx context.Context, domain string, localPart string) (*Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/domains/%s/aliases/%s", c.Endpoint, ascii, localPart)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	response := Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
