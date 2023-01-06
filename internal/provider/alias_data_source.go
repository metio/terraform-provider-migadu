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
	_ datasource.DataSource              = &aliasDataSource{}
	_ datasource.DataSourceWithConfigure = &aliasDataSource{}
)

func NewAliasDataSource() datasource.DataSource {
	return &aliasDataSource{}
}

type aliasDataSource struct {
	migaduClient *client.MigaduClient
}

type aliasDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
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

func (d *aliasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alias"
}

func (d *aliasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Get information about an email alias.",
		MarkdownDescription: "Get information about an email alias.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the alias.",
				MarkdownDescription: "The domain name of the alias.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the alias.",
				MarkdownDescription: "The local part of the alias.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Contains the value 'local_part@domain_name'.",
				MarkdownDescription: "Contains the value `local_part@domain_name`.",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				Description:         "The email address of the alias 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "The email address of the alias `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
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
				Description:         "Whether the alias is internal only. An internal alias can only receive emails from Migadu servers.",
				MarkdownDescription: "Whether the alias is internal only. An internal alias can only receive emails from Migadu servers.",
				Computed:            true,
			},
			"expirable": schema.BoolAttribute{
				Description:         "Whether the alias expires at some time.",
				MarkdownDescription: "Whether the alias expires at some time.",
				Computed:            true,
			},
			"expires_on": schema.StringAttribute{
				Description:         "The expiration date of the alias.",
				MarkdownDescription: "The expiration date of the alias.",
				Computed:            true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Description:         "Whether to remove the alias upon expiry.",
				MarkdownDescription: "Whether to remove the alias upon expiry.",
				Computed:            true,
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
			fmt.Sprintf("Expected *model.MigaduClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.migaduClient = migaduClient
}

func (d *aliasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data aliasDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alias, err := d.migaduClient.GetAlias(ctx, data.DomainName.ValueString(), data.LocalPart.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

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

	data.ID = types.StringValue(fmt.Sprintf("%s@%s", data.LocalPart.ValueString(), data.DomainName.ValueString()))
	data.Address = types.StringValue(alias.Address)
	data.Destinations = destinations
	data.DestinationsPunycode = destinationsPunycode
	data.IsInternal = types.BoolValue(alias.IsInternal)
	data.Expirable = types.BoolValue(alias.Expirable)
	data.ExpiresOn = types.StringValue(alias.ExpiresOn)
	data.RemoveUponExpiry = types.BoolValue(alias.RemoveUponExpiry)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
