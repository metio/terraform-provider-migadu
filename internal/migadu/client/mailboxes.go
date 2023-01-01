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
	"github.com/metio/terraform-provider-migadu/internal/idn"
	"github.com/metio/terraform-provider-migadu/internal/migadu/model"
	"golang.org/x/net/idna"
	"net/http"
)

// GetMailboxes - Returns mailboxes for a single domain
func (c *MigaduClient) GetMailboxes(ctx context.Context, domain string) (*model.Mailboxes, error) {
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

	response := model.Mailboxes{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetMailboxes: %w", err)
	}

	return &response, nil
}

// GetMailbox - Returns specific mailbox
func (c *MigaduClient) GetMailbox(ctx context.Context, domain string, localPart string) (*model.Mailbox, error) {
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

	response := model.Mailbox{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetMailbox: %w", err)
	}

	return &response, nil
}

// CreateMailbox - Creates a new mailbox
func (c *MigaduClient) CreateMailbox(ctx context.Context, domain string, mailbox *model.Mailbox) (*model.Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/mailboxes", c.Endpoint, ascii)

	senderDenyListASCII, err := idn.ConvertEmailsToASCII(mailbox.SenderDenyList)
	if err != nil {
		return nil, err
	}
	mailbox.SenderDenyList = senderDenyListASCII
	senderAllowListASCII, err := idn.ConvertEmailsToASCII(mailbox.SenderAllowList)
	if err != nil {
		return nil, err
	}
	mailbox.SenderAllowList = senderAllowListASCII
	recipientDenyListASCII, err := idn.ConvertEmailsToASCII(mailbox.RecipientDenyList)
	if err != nil {
		return nil, err
	}
	mailbox.RecipientDenyList = recipientDenyListASCII

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

	response := model.Mailbox{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}

	return &response, nil
}

// UpdateMailbox - Updates an existing mailbox
func (c *MigaduClient) UpdateMailbox(ctx context.Context, domain string, localPart string, mailbox *model.Mailbox) (*model.Mailbox, error) {
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

	response := model.Mailbox{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}

	return &response, nil
}

// DeleteMailbox - Deletes an existing mailbox
func (c *MigaduClient) DeleteMailbox(ctx context.Context, domain string, localPart string) (*model.Mailbox, error) {
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

	response := model.Mailbox{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("DeleteMailbox: %w", err)
	}

	return &response, nil
}
