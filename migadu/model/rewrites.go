/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package model

type Rewrites struct {
	Rewrites []Rewrite `json:"rewrites"`
}

type Rewrite struct {
	DomainName    string   `json:"domain_name"`
	Name          string   `json:"name"`
	LocalPartRule string   `json:"local_part_rule"`
	OrderNum      int64    `json:"order_num"`
	Destinations  []string `json:"destinations"`
}
