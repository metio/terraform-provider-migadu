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

type Mailboxes struct {
	Mailboxes []Mailbox `json:"mailboxes"`
}

type Mailbox struct {
	LocalPart             string   `json:"local_part"`
	DomainName            string   `json:"domain_name"`
	Address               string   `json:"address"`
	Name                  string   `json:"name"`
	IsInternal            bool     `json:"is_internal"`
	MaySend               bool     `json:"may_send"`
	MayReceive            bool     `json:"may_receive"`
	MayAccessImap         bool     `json:"may_access_imap"`
	MayAccessPop3         bool     `json:"may_access_pop3"`
	MayAccessManageSieve  bool     `json:"may_access_managesieve"`
	PasswordMethod        string   `json:"password_method"`
	Password              string   `json:"password"`
	PasswordRecoveryEmail string   `json:"password_recovery_email"`
	SpamAction            string   `json:"spam_action"`
	SpamAggressiveness    string   `json:"spam_aggressiveness"`
	Expirable             bool     `json:"expireable"`
	ExpiresOn             string   `json:"expires_on"`
	RemoveUponExpiry      bool     `json:"remove_upon_expiry"`
	SenderDenyList        []string `json:"sender_denylist"`
	SenderAllowList       []string `json:"sender_allowlist"`
	RecipientDenyList     []string `json:"recipient_denylist"`
	AutoRespondActive     bool     `json:"autorespond_active"`
	AutoRespondSubject    string   `json:"autorespond_subject"`
	AutoRespondBody       string   `json:"autorespond_body"`
	AutoRespondExpiresOn  string   `json:"autorespond_expires_on"`
	FooterActive          bool     `json:"footer_active"`
	FooterPlainBody       string   `json:"footer_plain_body"`
	FooterHtmlBody        string   `json:"footer_html_body"`
	StorageUsage          float64  `json:"storage_usage"`
	Delegations           []string `json:"delegations"`
	Identities            []string `json:"identities"`
}

// GetMailboxes - Returns mailboxes for a single domain
func (c *MigaduClient) GetMailboxes(domain string) (*Mailboxes, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/domains/%s/mailboxes", c.Endpoint, ascii)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	mailboxes := Mailboxes{}
	err = json.Unmarshal(body, &mailboxes)
	if err != nil {
		return nil, err
	}

	return &mailboxes, nil
}

// GetMailbox - Returns specific mailbox
func (c *MigaduClient) GetMailbox(domain string, localPart string) (*Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s", c.Endpoint, ascii, localPart)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	mailbox := Mailbox{}
	err = json.Unmarshal(body, &mailbox)
	if err != nil {
		return nil, err
	}

	return &mailbox, nil
}
