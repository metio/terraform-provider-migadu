/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package model

// Rewrites is the data model that wraps multiple rewrite rules
type Rewrites struct {
	Rewrites []Rewrite `json:"rewrites"`
}

// Rewrite is the data model for a single rewrite rule
type Rewrite struct {
	DomainName    string   `json:"domain_name"`
	Name          string   `json:"name"`
	LocalPartRule string   `json:"local_part_rule"`
	OrderNum      int64    `json:"order_num"`
	Destinations  []string `json:"destinations"`
}
