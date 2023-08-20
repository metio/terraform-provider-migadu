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

type rewriteRuleJson struct {
	*model.RewriteRule
	Destinations string `json:"destinations"`
}

// GetRewriteRules returns rewrite rules for a single domain
func (c *MigaduClient) GetRewriteRules(ctx context.Context, domain string) (*model.RewriteRules, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRules: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites", c.Endpoint, ascii)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRules: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRules: %w", err)
	}

	response := model.RewriteRules{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRules: %w", err)
	}

	return &response, nil
}

// GetRewriteRule returns a specific rewrite rule
func (c *MigaduClient) GetRewriteRule(ctx context.Context, domain string, slug string) (*model.RewriteRule, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRule: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites/%s", c.Endpoint, ascii, slug)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRule: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRule: %w", err)
	}

	response := model.RewriteRule{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRule: %w", err)
	}

	return &response, nil
}

// CreateRewriteRule creates a new rewrite rule
func (c *MigaduClient) CreateRewriteRule(ctx context.Context, domain string, rewrite *model.RewriteRule) (*model.RewriteRule, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateRewriteRule: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites", c.Endpoint, ascii)

	var requestBody []byte
	if rewrite != nil {
		asciiEmails, err := idn.ConvertEmailsToASCII(rewrite.Destinations)
		if err != nil {
			return nil, err
		}
		rewrite.Destinations = asciiEmails

		requestBody, err = json.Marshal(rewriteRuleJson{RewriteRule: rewrite, Destinations: strings.Join(rewrite.Destinations, ",")})
		if err != nil {
			return nil, fmt.Errorf("CreateRewriteRule: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("CreateRewriteRule: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("CreateRewriteRule: %w", err)
	}

	response := model.RewriteRule{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("CreateRewriteRule: %w", err)
	}

	return &response, nil
}

// UpdateRewriteRule updates an existing rewrite rule
func (c *MigaduClient) UpdateRewriteRule(ctx context.Context, domain string, slug string, rewrite *model.RewriteRule) (*model.RewriteRule, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateRewriteRule: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites/%s", c.Endpoint, ascii, slug)

	var requestBody []byte
	if rewrite != nil {
		asciiEmails, err := idn.ConvertEmailsToASCII(rewrite.Destinations)
		if err != nil {
			return nil, err
		}
		rewrite.Destinations = asciiEmails

		requestBody, err = json.Marshal(rewriteRuleJson{RewriteRule: rewrite, Destinations: strings.Join(rewrite.Destinations, ",")})
		if err != nil {
			return nil, fmt.Errorf("UpdateRewriteRule: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("UpdateRewriteRule: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("UpdateRewriteRule: %w", err)
	}

	response := model.RewriteRule{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("UpdateRewriteRule: %w", err)
	}

	return &response, nil
}

// DeleteRewriteRule deletes an existing rewrite rule
func (c *MigaduClient) DeleteRewriteRule(ctx context.Context, domain string, slug string) (*model.RewriteRule, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewriteRule: %w", err)
	}

	url := fmt.Sprintf("%s/domains/%s/rewrites/%s", c.Endpoint, ascii, slug)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewriteRule: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewriteRule: %w", err)
	}

	response := model.RewriteRule{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewriteRule: %w", err)
	}

	return &response, nil
}
