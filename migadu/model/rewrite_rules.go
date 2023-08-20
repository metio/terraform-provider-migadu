/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package model

// RewriteRules is the data model that wraps multiple rewrite rules
type RewriteRules struct {
	RewriteRules []RewriteRule `json:"rewrites"`
}

// RewriteRule is the data model for a single rewrite rule
type RewriteRule struct {
	DomainName    string   `json:"domain_name"`
	Name          string   `json:"name"`
	LocalPartRule string   `json:"local_part_rule"`
	OrderNum      int64    `json:"order_num"`
	Destinations  []string `json:"destinations"`
}
