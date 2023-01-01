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
	"github.com/metio/terraform-provider-migadu/internal/migadu/client"
)

var (
	_ datasource.DataSource              = &identityDataSource{}
	_ datasource.DataSourceWithConfigure = &identityDataSource{}
)

func NewIdentityDataSource() datasource.DataSource {
	return &identityDataSource{}
}

type identityDataSource struct {
	migaduClient *client.MigaduClient
}

type identityDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	LocalPart            types.String `tfsdk:"local_part"`
	DomainName           types.String `tfsdk:"domain_name"`
	Identity             types.String `tfsdk:"identity"`
	Address              types.String `tfsdk:"address"`
	Name                 types.String `tfsdk:"name"`
	MaySend              types.Bool   `tfsdk:"may_send"`
	MayReceive           types.Bool   `tfsdk:"may_receive"`
	MayAccessImap        types.Bool   `tfsdk:"may_access_imap"`
	MayAccessPop3        types.Bool   `tfsdk:"may_access_pop3"`
	MayAccessManageSieve types.Bool   `tfsdk:"may_access_manage_sieve"`
	FooterActive         types.Bool   `tfsdk:"footer_active"`
	FooterPlainBody      types.String `tfsdk:"footer_plain_body"`
	FooterHtmlBody       types.String `tfsdk:"footer_html_body"`
}

func (d *identityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity"
}

func (d *identityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gets a single identity.",
		MarkdownDescription: "Gets a single identity.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain to fetch identities of.",
				MarkdownDescription: "The domain to fetch identities of.",
				Required:            true,
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part to fetch identities of.",
				MarkdownDescription: "The local part to fetch identities of.",
				Required:            true,
			},
			"identity": schema.StringAttribute{
				Description:         "The local part of the identity to fetch.",
				MarkdownDescription: "The local part of the identity to fetch.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Contains the value 'local_part@domain_name/identity'.",
				MarkdownDescription: "Contains the value `local_part@domain_name/identity`.",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"may_send": schema.BoolAttribute{
				Computed: true,
			},
			"may_receive": schema.BoolAttribute{
				Computed: true,
			},
			"may_access_imap": schema.BoolAttribute{
				Computed: true,
			},
			"may_access_pop3": schema.BoolAttribute{
				Computed: true,
			},
			"may_access_manage_sieve": schema.BoolAttribute{
				Computed: true,
			},
			"footer_active": schema.BoolAttribute{
				Computed: true,
			},
			"footer_plain_body": schema.StringAttribute{
				Computed: true,
			},
			"footer_html_body": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *identityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *identityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data identityDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity, err := d.migaduClient.GetIdentity(ctx, data.DomainName.ValueString(), data.LocalPart.ValueString(), data.Identity.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

	//data.LocalPart = types.StringValue(identity.LocalPart)
	//data.DomainName = types.StringValue(identity.DomainName)
	//data.Identity = types.StringValue(identity.Identity)
	data.Address = types.StringValue(identity.Address)
	data.Name = types.StringValue(identity.Name)
	data.MaySend = types.BoolValue(identity.MaySend)
	data.MayReceive = types.BoolValue(identity.MayReceive)
	data.MayAccessImap = types.BoolValue(identity.MayAccessImap)
	data.MayAccessPop3 = types.BoolValue(identity.MayAccessPop3)
	data.MayAccessManageSieve = types.BoolValue(identity.MayAccessManageSieve)
	data.FooterActive = types.BoolValue(identity.FooterActive)
	data.FooterPlainBody = types.StringValue(identity.FooterPlainBody)
	data.FooterHtmlBody = types.StringValue(identity.FooterHtmlBody)

	data.ID = types.StringValue(fmt.Sprintf("%s@%s/%s", data.LocalPart.ValueString(), data.DomainName.ValueString(), data.Identity.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
