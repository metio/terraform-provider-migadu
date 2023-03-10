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
	"github.com/metio/terraform-provider-migadu/migadu/idn"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"golang.org/x/net/idna"
	"net/http"
	"strings"
)

type rewriteJson struct {
	*model.Rewrite
	Destinations string `json:"destinations"`
}

// GetRewrites returns rewrite rules for a single domain
func (c *MigaduClient) GetRewrites(ctx context.Context, domain string) (*model.Rewrites, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetRewrites: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites", c.Endpoint, ascii)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetRewrites: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetRewrites: %w", err)
	}

	response := model.Rewrites{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetRewrites: %w", err)
	}

	return &response, nil
}

// GetRewrite returns a specific rewrite rule
func (c *MigaduClient) GetRewrite(ctx context.Context, domain string, slug string) (*model.Rewrite, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetRewrite: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites/%s", c.Endpoint, ascii, slug)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetRewrite: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetRewrite: %w", err)
	}

	response := model.Rewrite{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetRewrite: %w", err)
	}

	return &response, nil
}

// CreateRewrite creates a new rewrite rule
func (c *MigaduClient) CreateRewrite(ctx context.Context, domain string, rewrite *model.Rewrite) (*model.Rewrite, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateRewrite: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites", c.Endpoint, ascii)

	var requestBody []byte
	if rewrite != nil {
		asciiEmails, err := idn.ConvertEmailsToASCII(rewrite.Destinations)
		if err != nil {
			return nil, err
		}
		rewrite.Destinations = asciiEmails

		requestBody, err = json.Marshal(rewriteJson{Rewrite: rewrite, Destinations: strings.Join(rewrite.Destinations, ",")})
		if err != nil {
			return nil, fmt.Errorf("CreateRewrite: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("CreateRewrite: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("CreateRewrite: %w", err)
	}

	response := model.Rewrite{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("CreateRewrite: %w", err)
	}

	return &response, nil
}

// UpdateRewrite updates an existing rewrite rule
func (c *MigaduClient) UpdateRewrite(ctx context.Context, domain string, slug string, rewrite *model.Rewrite) (*model.Rewrite, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateRewrite: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites/%s", c.Endpoint, ascii, slug)

	var requestBody []byte
	if rewrite != nil {
		asciiEmails, err := idn.ConvertEmailsToASCII(rewrite.Destinations)
		if err != nil {
			return nil, err
		}
		rewrite.Destinations = asciiEmails

		requestBody, err = json.Marshal(rewriteJson{Rewrite: rewrite, Destinations: strings.Join(rewrite.Destinations, ",")})
		if err != nil {
			return nil, fmt.Errorf("UpdateRewrite: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("UpdateRewrite: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("UpdateRewrite: %w", err)
	}

	response := model.Rewrite{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("UpdateRewrite: %w", err)
	}

	return &response, nil
}

// DeleteRewrite deletes an existing rewrite rule
func (c *MigaduClient) DeleteRewrite(ctx context.Context, domain string, slug string) (*model.Rewrite, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewrite: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites/%s", c.Endpoint, ascii, slug)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewrite: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewrite: %w", err)
	}

	response := model.Rewrite{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewrite: %w", err)
	}

	return &response, nil
}
