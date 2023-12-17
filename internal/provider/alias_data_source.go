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
	_ datasource.DataSource              = (*AliasDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*AliasDataSource)(nil)
)

func NewAliasDataSource() datasource.DataSource {
	return &AliasDataSource{}
}

type AliasDataSource struct {
	migaduClient *client.MigaduClient
}

type AliasDataSourceModel struct {
	ID               custom_types.EmailAddressValue    `tfsdk:"id"`
	LocalPart        types.String                      `tfsdk:"local_part"`
	DomainName       custom_types.DomainNameValue      `tfsdk:"domain_name"`
	Address          custom_types.EmailAddressValue    `tfsdk:"address"`
	Destinations     custom_types.EmailAddressSetValue `tfsdk:"destinations"`
	IsInternal       types.Bool                        `tfsdk:"is_internal"`
	Expirable        types.Bool                        `tfsdk:"expirable"`
	ExpiresOn        types.String                      `tfsdk:"expires_on"`
	RemoveUponExpiry types.Bool                        `tfsdk:"remove_upon_expiry"`
}

func (d *AliasDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_alias"
}

func (d *AliasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Get information about an email alias.",
		MarkdownDescription: "Get information about an email alias.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Contains the value 'local_part@domain_name'.",
				MarkdownDescription: "Contains the value `local_part@domain_name`.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the alias.",
				MarkdownDescription: "The local part of the alias.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the alias.",
				MarkdownDescription: "The domain name of the alias.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType:          custom_types.DomainNameType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"address": schema.StringAttribute{
				Description:         "The email address of the alias 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "The email address of the alias `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
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
				Description:         "Whether the alias is internal only. An internal alias can only receive emails from Migadu servers.",
				MarkdownDescription: "Whether the alias is internal only. An internal alias can only receive emails from Migadu servers.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"expirable": schema.BoolAttribute{
				Description:         "Whether the alias expires at some time.",
				MarkdownDescription: "Whether the alias expires at some time.",
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
				Description:         "Whether to remove the alias upon expiry.",
				MarkdownDescription: "Whether to remove the alias upon expiry.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
		},
	}
}

func (d *AliasDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if migaduClient, ok := request.ProviderData.(*client.MigaduClient); ok {
		d.migaduClient = migaduClient
	} else {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *model.MigaduClient, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
}

func (d *AliasDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data AliasDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	alias, err := d.migaduClient.GetAlias(ctx, data.DomainName.ValueString(), data.LocalPart.ValueString())
	if err != nil {
		response.Diagnostics.Append(AliasReadError(err))
		return
	}

	destinations, diags := custom_types.NewEmailAddressSetValueFrom(ctx, alias.Destinations)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data.ID = custom_types.NewEmailAddressValue(CreateAliasID(data.LocalPart, data.DomainName))
	data.Address = custom_types.NewEmailAddressValue(alias.Address)
	data.Destinations = destinations
	data.IsInternal = types.BoolValue(alias.IsInternal)
	data.Expirable = types.BoolValue(alias.Expirable)
	data.ExpiresOn = types.StringValue(alias.ExpiresOn)
	data.RemoveUponExpiry = types.BoolValue(alias.RemoveUponExpiry)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
