/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/metio/terraform-provider-migadu/internal/client"
)

var (
	_ datasource.DataSource              = &aliasDataSource{}
	_ datasource.DataSourceWithConfigure = &aliasDataSource{}
)

func NewAliasDataSource() datasource.DataSource {
	return &aliasDataSource{}
}

type aliasDataSource struct {
	migaduClient *client.MigaduClient
}

type AliasDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	LocalPart        types.String `tfsdk:"local_part"`
	DomainName       types.String `tfsdk:"domain_name"`
	Address          types.String `tfsdk:"address"`
	Destinations     types.List   `tfsdk:"destinations"`
	IsInternal       types.Bool   `tfsdk:"is_internal"`
	Expirable        types.Bool   `tfsdk:"expirable"`
	ExpiresOn        types.String `tfsdk:"expires_on"`
	RemoveUponExpiry types.Bool   `tfsdk:"remove_upon_expiry"`
}

func (d *aliasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alias"
}

func (d *aliasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gets a single alias.",
		MarkdownDescription: "Gets a single alias.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the alias to fetch.",
				MarkdownDescription: "The domain name of the alias to fetch.",
				Required:            true,
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the alias to fetch.",
				MarkdownDescription: "The local part of the alias to fetch.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Contains the full email address 'local_part@domain_name'.",
				MarkdownDescription: "Contains the full email address `local_part@domain_name`.",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				Computed: true,
			},
			"destinations": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_internal": schema.BoolAttribute{
				Computed: true,
			},
			"expirable": schema.BoolAttribute{
				Computed: true,
			},
			"expires_on": schema.StringAttribute{
				Computed: true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *aliasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	migaduClient, ok := req.ProviderData.(*client.MigaduClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.MigaduClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.migaduClient = migaduClient
}

func (d *aliasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AliasDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alias, err := d.migaduClient.GetAlias(ctx, data.DomainName.ValueString(), data.LocalPart.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

	data.Address = types.StringValue(alias.Address)
	data.IsInternal = types.BoolValue(alias.IsInternal)
	data.Expirable = types.BoolValue(alias.Expirable)
	data.ExpiresOn = types.StringValue(alias.ExpiresOn)
	data.RemoveUponExpiry = types.BoolValue(alias.RemoveUponExpiry)

	destinations, diags := types.ListValueFrom(ctx, types.StringType, alias.Destinations)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Destinations = destinations

	data.ID = types.StringValue(fmt.Sprintf("%s@%s", data.LocalPart.ValueString(), data.DomainName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
