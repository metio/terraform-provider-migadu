/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/metio/terraform-provider-migadu/migadu/client"
)

var (
	_ datasource.DataSource              = &rewriteDataSource{}
	_ datasource.DataSourceWithConfigure = &rewriteDataSource{}
)

func NewRewriteDataSource() datasource.DataSource {
	return &rewriteDataSource{}
}

type rewriteDataSource struct {
	migaduClient *client.MigaduClient
}

type rewriteDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	DomainName           types.String `tfsdk:"domain_name"`
	Name                 types.String `tfsdk:"name"`
	LocalPartRule        types.String `tfsdk:"local_part_rule"`
	OrderNum             types.Int64  `tfsdk:"order_num"`
	Destinations         types.List   `tfsdk:"destinations"`
	DestinationsPunycode types.List   `tfsdk:"destinations_punycode"`
}

func (d *rewriteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rewrite"
}

func (d *rewriteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Get information about a single rewrite rule.",
		MarkdownDescription: "Get information about a single rewrite rule.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain of the rewrite rule.",
				MarkdownDescription: "The domain of the rewrite rule.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name (slug) of the rewrite rule.",
				MarkdownDescription: "The name (slug) of the rewrite rule.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Contains the value 'domain_name/name'.",
				MarkdownDescription: "Contains the value `domain_name/name`.",
				Computed:            true,
			},
			"local_part_rule": schema.StringAttribute{
				Description:         "The local part expression of the rewrite rule",
				MarkdownDescription: "The local part expression of the rewrite rule",
				Computed:            true,
			},
			"order_num": schema.Int64Attribute{
				Description:         "The order number of the rewrite rule.",
				MarkdownDescription: "The order number of the rewrite rule.",
				Computed:            true,
			},
			"destinations": schema.ListAttribute{
				Description:         "The destinations of the rewrite rule in unicode.",
				MarkdownDescription: "The destinations of the rewrite rule in unicode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"destinations_punycode": schema.ListAttribute{
				Description:         "The destinations of the rewrite rule in punycode.",
				MarkdownDescription: "The destinations of the rewrite rule in punycode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *rewriteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *rewriteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data rewriteDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rewrite, err := d.migaduClient.GetRewrite(ctx, data.DomainName.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

	destinations, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(rewrite.Destinations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	destinationsPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(rewrite.Destinations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Destinations = destinations
	data.DestinationsPunycode = destinationsPunycode
	data.ID = types.StringValue(fmt.Sprintf("%s/%s", data.DomainName.ValueString(), data.Name.ValueString()))
	data.LocalPartRule = types.StringValue(rewrite.LocalPartRule)
	data.OrderNum = types.Int64Value(rewrite.OrderNum)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
