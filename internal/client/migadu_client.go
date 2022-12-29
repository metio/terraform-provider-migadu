/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MigaduClient struct {
	HTTPClient *http.Client
	Endpoint   string
	Username   string
	Token      string
}

func New(endpoint, username, token *string, timeout time.Duration) (*MigaduClient, error) {
	if username == nil {
		return nil, errors.New("no username supplied")
	}
	if token == nil {
		return nil, errors.New("no token supplied")
	}

	c := MigaduClient{
		HTTPClient: &http.Client{Timeout: timeout},
		Endpoint:   "https://api.migadu.com/v1/",
		Username:   *username,
		Token:      *token,
	}

	if endpoint != nil {
		c.Endpoint = *endpoint
	}

	return &c, nil
}

func (c *MigaduClient) doRequest(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.Username, c.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
