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
	"github.com/metio/migadu-client.go/client"
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
)

var (
	_ datasource.DataSource              = (*RewriteRulesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*RewriteRulesDataSource)(nil)
)

func NewRewriteRulesDataSource() datasource.DataSource {
	return &RewriteRulesDataSource{}
}

type RewriteRulesDataSource struct {
	MigaduClient *client.MigaduClient
}

type RewriteRulesDataSourceModel struct {
	ID         custom_types.DomainNameValue `tfsdk:"id"`
	DomainName custom_types.DomainNameValue `tfsdk:"domain_name"`
	Rewrites   []RewriteRuleModel           `tfsdk:"rewrites"`
}

type RewriteRuleModel struct {
	DomainName    custom_types.DomainNameValue      `tfsdk:"domain_name"`
	Name          types.String                      `tfsdk:"name"`
	LocalPartRule types.String                      `tfsdk:"local_part_rule"`
	OrderNum      types.Int64                       `tfsdk:"order_num"`
	Destinations  custom_types.EmailAddressSetValue `tfsdk:"destinations"`
}

func (d *RewriteRulesDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rewrite_rules"
}

func (d *RewriteRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Get information about a all rewrite rules of a domain.",
		MarkdownDescription: "Get information about a all rewrite rules of a domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Same value as the 'domain_name' attribute.",
				MarkdownDescription: "Same value as the `domain_name` attribute.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.DomainNameType{},
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain to fetch rewrite rules of.",
				MarkdownDescription: "The domain to fetch rewrite rules of.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType:          custom_types.DomainNameType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"rewrites": schema.ListNestedAttribute{
				Description:         "The configured rewrite rules for the given 'domain_name'.",
				MarkdownDescription: "The configured rewrite rules for the given `domain_name`.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"local_part_rule": schema.StringAttribute{
							Description:         "The local part expression of the rewrite rule",
							MarkdownDescription: "The local part expression of the rewrite rule",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"domain_name": schema.StringAttribute{
							Description:         "The domain of the rewrite rule.",
							MarkdownDescription: "The domain of the rewrite rule.",
							Required:            false,
							Optional:            false,
							Computed:            true,
							CustomType:          custom_types.DomainNameType{},
						},
						"name": schema.StringAttribute{
							Description:         "The name (slug) of the rewrite rule.",
							MarkdownDescription: "The name (slug) of the rewrite rule.",
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
				},
			},
		},
	}
}

func (d *RewriteRulesDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (d *RewriteRulesDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data RewriteRulesDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	rewrites, err := d.MigaduClient.GetRewriteRules(ctx, data.DomainName.ValueString())
	if err != nil {
		response.Diagnostics.Append(RewriteRuleReadError(err))
		return
	}

	for _, rewrite := range rewrites.RewriteRules {
		destinations, diags := custom_types.NewEmailAddressSetValueFrom(ctx, rewrite.Destinations)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		model := RewriteRuleModel{
			DomainName:    custom_types.NewDomainNameValue(rewrite.DomainName),
			Name:          types.StringValue(rewrite.Name),
			LocalPartRule: types.StringValue(rewrite.LocalPartRule),
			OrderNum:      types.Int64Value(rewrite.OrderNum),
			Destinations:  destinations,
		}

		data.Rewrites = append(data.Rewrites, model)
	}

	data.ID = data.DomainName

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
