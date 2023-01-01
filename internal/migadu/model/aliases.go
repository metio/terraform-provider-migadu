/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package model

type Aliases struct {
	Aliases []Alias `json:"address_aliases"`
}

type Alias struct {
	LocalPart        string   `json:"local_part"`
	DomainName       string   `json:"domain_name"`
	Address          string   `json:"address"`
	Destinations     []string `json:"destinations"`
	IsInternal       bool     `json:"is_internal"`
	Expirable        bool     `json:"expireable"`
	ExpiresOn        string   `json:"expires_on"`
	RemoveUponExpiry bool     `json:"remove_upon_expiry"`
}
