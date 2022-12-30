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
	_ datasource.DataSource              = &RewritesDataSource{}
	_ datasource.DataSourceWithConfigure = &RewritesDataSource{}
)

func NewRewritesDataSource() datasource.DataSource {
	return &RewritesDataSource{}
}

type RewritesDataSource struct {
	migaduClient *client.MigaduClient
}

type RewritesDataSourceModel struct {
	ID         types.String   `tfsdk:"id"`
	DomainName types.String   `tfsdk:"domain_name"`
	Rewrites   []RewriteModel `tfsdk:"rewrites"`
}

type RewriteModel struct {
	DomainName    types.String `tfsdk:"domain_name"`
	Name          types.String `tfsdk:"name"`
	LocalPartRule types.String `tfsdk:"local_part_rule"`
	OrderNum      types.Int64  `tfsdk:"order_num"`
	Destinations  types.List   `tfsdk:"destinations"`
}

func (d *RewritesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rewrites"
}

func (d *RewritesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gets all rewrites of a domain.",
		MarkdownDescription: "Gets all rewrites of a domain.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain to fetch rewrites of.",
				MarkdownDescription: "The domain to fetch rewrites of.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Same value as the 'domain_name' attribute.",
				MarkdownDescription: "Same value as the `domain_name` attribute.",
				Computed:            true,
			},
			"rewrites": schema.ListNestedAttribute{
				Description:         "The configured rewrites for the given 'domain_name'.",
				MarkdownDescription: "The configured rewrites for the given `domain_name`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain_name": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"local_part_rule": schema.StringAttribute{
							Computed: true,
						},
						"order_num": schema.Int64Attribute{
							Computed: true,
						},
						"destinations": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *RewritesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RewritesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RewritesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rewrites, err := d.migaduClient.GetRewrites(ctx, data.DomainName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

	for _, rewrite := range rewrites.Rewrites {
		aliasModel := RewriteModel{
			DomainName:    types.StringValue(rewrite.DomainName),
			Name:          types.StringValue(rewrite.Name),
			LocalPartRule: types.StringValue(rewrite.LocalPartRule),
			OrderNum:      types.Int64Value(rewrite.OrderNum),
		}

		destinations, diags := types.ListValueFrom(ctx, types.StringType, rewrite.Destinations)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		aliasModel.Destinations = destinations

		data.Rewrites = append(data.Rewrites, aliasModel)
	}

	data.ID = data.DomainName

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
