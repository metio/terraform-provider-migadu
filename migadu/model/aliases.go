/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package model

// Aliases is the data model that wraps multiple aliases
type Aliases struct {
	Aliases []Alias `json:"address_aliases"`
}

// Alias is the data model for a single alias
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
