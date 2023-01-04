/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package client_test

import (
	"github.com/metio/terraform-provider-migadu/migadu/client"
	"time"
)

func newTestClient(endpoint string) *client.MigaduClient {
	username := "username"
	token := "token"
	c, err := client.New(&endpoint, &username, &token, 10*time.Second)
	if err != nil {
		panic(err)
	}
	return c
}
