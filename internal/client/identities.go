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
func (c *MigaduClient) GetIdentities(domain string, localPart string) (*Identities, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s/identities", c.Endpoint, ascii, localPart)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	identities := Identities{}
	err = json.Unmarshal(body, &identities)
	if err != nil {
		return nil, err
	}

	return &identities, nil
}

// GetIdentity - Returns specific identity
func (c *MigaduClient) GetIdentity(domain string, localPart string, id string) (*Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s/identities/%s", c.Endpoint, ascii, localPart, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	identity := Identity{}
	err = json.Unmarshal(body, &identity)
	if err != nil {
		return nil, err
	}

	return &identity, nil
}
