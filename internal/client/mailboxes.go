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
func (c *MigaduClient) GetMailboxes(ctx context.Context, domain string) (*Mailboxes, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetMailboxes: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes", c.Endpoint, ascii)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetMailboxes: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetMailboxes: %w", err)
	}

	response := Mailboxes{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetMailboxes: %w", err)
	}

	return &response, nil
}

// GetMailbox - Returns specific mailbox
func (c *MigaduClient) GetMailbox(ctx context.Context, domain string, localPart string) (*Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetMailbox: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s", c.Endpoint, ascii, localPart)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetMailbox: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetMailbox: %w", err)
	}

	response := Mailbox{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetMailbox: %w", err)
	}

	return &response, nil
}

// CreateMailbox - Creates a new mailbox
func (c *MigaduClient) CreateMailbox(ctx context.Context, domain string, mailbox *Mailbox) (*Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes", c.Endpoint, ascii)

	requestBody, err := json.Marshal(mailbox)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}

	response := Mailbox{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}

	return &response, nil
}

// UpdateMailbox - Updates an existing mailbox
func (c *MigaduClient) UpdateMailbox(ctx context.Context, domain string, localPart string, mailbox *Mailbox) (*Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s", c.Endpoint, ascii, localPart)

	requestBody, err := json.Marshal(mailbox)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	response := Mailbox{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	return &response, nil
}

// DeleteMailbox - Deletes an existing mailbox
func (c *MigaduClient) DeleteMailbox(ctx context.Context, domain string, localPart string) (*Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteMailbox: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes/%s", c.Endpoint, ascii, localPart)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("DeleteMailbox: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("DeleteMailbox: %w", err)
	}

	response := Mailbox{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("DeleteMailbox: %w", err)
	}

	return &response, nil
}
