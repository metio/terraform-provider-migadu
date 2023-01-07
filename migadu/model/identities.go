/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package model

// Identities is the data model that wraps multiple identities
type Identities struct {
	Identities []Identity `json:"identities"`
}

// Identity is the data model for a single identity
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
