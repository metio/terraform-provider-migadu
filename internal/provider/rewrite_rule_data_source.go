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
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
	"github.com/metio/terraform-provider-migadu/migadu/client"
)

var (
	_ datasource.DataSource              = (*RewriteRuleDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*RewriteRuleDataSource)(nil)
)

func NewRewriteRuleDataSource() datasource.DataSource {
	return &RewriteRuleDataSource{}
}

type RewriteRuleDataSource struct {
	MigaduClient *client.MigaduClient
}

type RewriteRuleDataSourceModel struct {
	ID            types.String                      `tfsdk:"id"`
	DomainName    custom_types.DomainNameValue      `tfsdk:"domain_name"`
	Name          types.String                      `tfsdk:"name"`
	LocalPartRule types.String                      `tfsdk:"local_part_rule"`
	OrderNum      types.Int64                       `tfsdk:"order_num"`
	Destinations  custom_types.EmailAddressSetValue `tfsdk:"destinations"`
}

func (d *RewriteRuleDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rewrite_rule"
}

func (d *RewriteRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Get information about a single rewrite rule.",
		MarkdownDescription: "Get information about a single rewrite rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Contains the value 'domain_name/name'.",
				MarkdownDescription: "Contains the value `domain_name/name`.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain of the rewrite rule.",
				MarkdownDescription: "The domain of the rewrite rule.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType:          custom_types.DomainNameType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name (slug) of the rewrite rule.",
				MarkdownDescription: "The name (slug) of the rewrite rule.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"local_part_rule": schema.StringAttribute{
				Description:         "The local part expression of the rewrite rule",
				MarkdownDescription: "The local part expression of the rewrite rule",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"order_num": schema.Int64Attribute{
				Description:         "The order number of the rewrite rule.",
				MarkdownDescription: "The order number of the rewrite rule.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"destinations": schema.SetAttribute{
				Description:         "The destinations of the rewrite rule.",
				MarkdownDescription: "The destinations of the rewrite rule.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType: custom_types.EmailAddressSetType{
					SetType: types.SetType{
						ElemType: custom_types.EmailAddressType{},
					},
				},
			},
		},
	}
}

func (d *RewriteRuleDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if migaduClient, ok := request.ProviderData.(*client.MigaduClient); ok {
		d.MigaduClient = migaduClient
	} else {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.MigaduClient, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
}

func (d *RewriteRuleDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data RewriteRuleDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	rewrite, err := d.MigaduClient.GetRewriteRule(ctx, data.DomainName.ValueString(), data.Name.ValueString())
	if err != nil {
		response.Diagnostics.Append(RewriteRuleReadError(err))
		return
	}

	destinations, diags := custom_types.NewEmailAddressSetValueFrom(ctx, rewrite.Destinations)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data.Destinations = destinations
	data.ID = types.StringValue(fmt.Sprintf("%s/%s", data.DomainName.ValueString(), data.Name.ValueString()))
	data.LocalPartRule = types.StringValue(rewrite.LocalPartRule)
	data.OrderNum = types.Int64Value(rewrite.OrderNum)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
