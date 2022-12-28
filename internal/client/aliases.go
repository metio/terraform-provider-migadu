/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
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
	req, err := http.NewRequest("GET", url, nil)
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
	req, err := http.NewRequest("GET", url, nil)
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
