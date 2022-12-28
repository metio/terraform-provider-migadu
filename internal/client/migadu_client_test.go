/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import "time"

func newTestClient(endpoint string) *MigaduClient {
	username := "username"
	token := "token"
	c, err := New(&endpoint, &username, &token, 10*time.Second)
	if err != nil {
		panic(err)
	}
	return c
}
