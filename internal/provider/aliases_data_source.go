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
	_ datasource.DataSource              = &aliasesDataSource{}
	_ datasource.DataSourceWithConfigure = &aliasesDataSource{}
)

func NewAliasesDataSource() datasource.DataSource {
	return &aliasesDataSource{}
}

type aliasesDataSource struct {
	migaduClient *client.MigaduClient
}

type aliasesDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	DomainName types.String `tfsdk:"domain_name"`
	Aliases    []aliasModel `tfsdk:"aliases"`
}

type aliasModel struct {
	LocalPart            types.String `tfsdk:"local_part"`
	DomainName           types.String `tfsdk:"domain_name"`
	Address              types.String `tfsdk:"address"`
	Destinations         types.List   `tfsdk:"destinations"`
	DestinationsPunycode types.List   `tfsdk:"destinations_punycode"`
	IsInternal           types.Bool   `tfsdk:"is_internal"`
	Expirable            types.Bool   `tfsdk:"expirable"`
	ExpiresOn            types.String `tfsdk:"expires_on"`
	RemoveUponExpiry     types.Bool   `tfsdk:"remove_upon_expiry"`
}

func (d *aliasesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aliases"
}

func (d *aliasesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Get information about all email aliases of a domain.",
		MarkdownDescription: "Get information about all email aliases of a domain.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of all aliases.",
				MarkdownDescription: "The domain name of all aliases.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Same value as the 'domain_name' attribute.",
				MarkdownDescription: "Same value as the `domain_name` attribute.",
				Computed:            true,
			},
			"aliases": schema.ListNestedAttribute{
				Description:         "The configured aliases for the given 'domain_name'.",
				MarkdownDescription: "The configured aliases for the given `domain_name`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain_name": schema.StringAttribute{
							Description:         "The domain name of the alias.",
							MarkdownDescription: "The domain name of the alias.",
							Computed:            true,
						},
						"local_part": schema.StringAttribute{
							Description:         "The local part of the alias.",
							MarkdownDescription: "The local part of the alias.",
							Computed:            true,
						},
						"address": schema.StringAttribute{
							Description:         "The email address 'local_part@domain_name' as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
							MarkdownDescription: "The email address `local_part@domain_name` as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
							Computed:            true,
						},
						"destinations": schema.ListAttribute{
							Description:         "List of email addresses that act as destinations of the alias in unicode.",
							MarkdownDescription: "List of email addresses that act as destinations of the alias in unicode.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"destinations_punycode": schema.ListAttribute{
							Description:         "List of email addresses that act as destinations of the alias in punycode.",
							MarkdownDescription: "List of email addresses that act as destinations of the alias in punycode.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"is_internal": schema.BoolAttribute{
							Description:         "Whether the alias is internal and can only receive emails from Migadu servers.",
							MarkdownDescription: "Whether the alias is internal and can only receive emails from Migadu servers.",
							Computed:            true,
						},
						"expirable": schema.BoolAttribute{
							Description:         "Whether the alias expires some time in the future.",
							MarkdownDescription: "Whether the alias expires some time in the future.",
							Computed:            true,
						},
						"expires_on": schema.StringAttribute{
							Description:         "The expiration date of the alias.",
							MarkdownDescription: "The expiration date of the alias.",
							Computed:            true,
						},
						"remove_upon_expiry": schema.BoolAttribute{
							Description:         "Whether the alias is removed once it is expired.",
							MarkdownDescription: "Whether the alias is removed once it is expired.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *aliasesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *aliasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data aliasesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	aliases, err := d.migaduClient.GetAliases(ctx, data.DomainName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

	for _, alias := range aliases.Aliases {
		destinations, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(alias.Destinations, &resp.Diagnostics))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		destinationsPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(alias.Destinations, &resp.Diagnostics))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		model := aliasModel{
			LocalPart:            types.StringValue(alias.LocalPart),
			DomainName:           types.StringValue(alias.DomainName),
			Destinations:         destinations,
			DestinationsPunycode: destinationsPunycode,
			Address:              types.StringValue(alias.Address),
			IsInternal:           types.BoolValue(alias.IsInternal),
			Expirable:            types.BoolValue(alias.Expirable),
			ExpiresOn:            types.StringValue(alias.ExpiresOn),
			RemoveUponExpiry:     types.BoolValue(alias.RemoveUponExpiry),
		}

		data.Aliases = append(data.Aliases, model)
	}

	data.ID = data.DomainName

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
