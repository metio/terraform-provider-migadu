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
	_ datasource.DataSource              = (*AliasesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*AliasesDataSource)(nil)
)

func NewAliasesDataSource() datasource.DataSource {
	return &AliasesDataSource{}
}

type AliasesDataSource struct {
	MigaduClient *client.MigaduClient
}

type AliasesDataSourceModel struct {
	ID         custom_types.DomainNameValue `tfsdk:"id"`
	DomainName custom_types.DomainNameValue `tfsdk:"domain_name"`
	Aliases    []AliasModel                 `tfsdk:"aliases"`
}

type AliasModel struct {
	LocalPart        types.String                      `tfsdk:"local_part"`
	DomainName       custom_types.DomainNameValue      `tfsdk:"domain_name"`
	Address          custom_types.EmailAddressValue    `tfsdk:"address"`
	Destinations     custom_types.EmailAddressSetValue `tfsdk:"destinations"`
	IsInternal       types.Bool                        `tfsdk:"is_internal"`
	Expirable        types.Bool                        `tfsdk:"expirable"`
	ExpiresOn        types.String                      `tfsdk:"expires_on"`
	RemoveUponExpiry types.Bool                        `tfsdk:"remove_upon_expiry"`
}

func (d *AliasesDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_aliases"
}

func (d *AliasesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Get information about all email aliases of a domain.",
		MarkdownDescription: "Get information about all email aliases of a domain.",
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
				Description:         "The domain name of all aliases.",
				MarkdownDescription: "The domain name of all aliases.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType:          custom_types.DomainNameType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"aliases": schema.ListNestedAttribute{
				Description:         "The configured aliases for the given 'domain_name'.",
				MarkdownDescription: "The configured aliases for the given `domain_name`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"local_part": schema.StringAttribute{
							Description:         "The local part of the alias.",
							MarkdownDescription: "The local part of the alias.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"domain_name": schema.StringAttribute{
							Description:         "The domain name of the alias.",
							MarkdownDescription: "The domain name of the alias.",
							Required:            false,
							Optional:            false,
							Computed:            true,
							CustomType:          custom_types.DomainNameType{},
						},
						"address": schema.StringAttribute{
							Description:         "The email address 'local_part@domain_name' as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
							MarkdownDescription: "The email address `local_part@domain_name` as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
							Required:            false,
							Optional:            false,
							Computed:            true,
							CustomType:          custom_types.EmailAddressType{},
						},
						"destinations": schema.SetAttribute{
							Description:         "List of email addresses that act as destinations of the alias.",
							MarkdownDescription: "List of email addresses that act as destinations of the alias.",
							Required:            false,
							Optional:            false,
							Computed:            true,
							CustomType: custom_types.EmailAddressSetType{
								SetType: types.SetType{
									ElemType: custom_types.EmailAddressType{},
								},
							},
						},
						"is_internal": schema.BoolAttribute{
							Description:         "Whether the alias is internal and can only receive emails from Migadu servers.",
							MarkdownDescription: "Whether the alias is internal and can only receive emails from Migadu servers.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"expirable": schema.BoolAttribute{
							Description:         "Whether the alias expires some time in the future.",
							MarkdownDescription: "Whether the alias expires some time in the future.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"expires_on": schema.StringAttribute{
							Description:         "The expiration date of the alias.",
							MarkdownDescription: "The expiration date of the alias.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"remove_upon_expiry": schema.BoolAttribute{
							Description:         "Whether the alias is removed once it is expired.",
							MarkdownDescription: "Whether the alias is removed once it is expired.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AliasesDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (d *AliasesDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data AliasesDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	aliases, err := d.MigaduClient.GetAliases(ctx, data.DomainName.ValueString())
	if err != nil {
		response.Diagnostics.Append(AliasReadError(err))
		return
	}

	for _, alias := range aliases.Aliases {
		destinations, diags := custom_types.NewEmailAddressSetValueFrom(ctx, alias.Destinations)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		model := AliasModel{
			LocalPart:        types.StringValue(alias.LocalPart),
			DomainName:       custom_types.NewDomainNameValue(alias.DomainName),
			Destinations:     destinations,
			Address:          custom_types.NewEmailAddressValue(alias.Address),
			IsInternal:       types.BoolValue(alias.IsInternal),
			Expirable:        types.BoolValue(alias.Expirable),
			ExpiresOn:        types.StringValue(alias.ExpiresOn),
			RemoveUponExpiry: types.BoolValue(alias.RemoveUponExpiry),
		}

		data.Aliases = append(data.Aliases, model)
	}

	data.ID = data.DomainName

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
