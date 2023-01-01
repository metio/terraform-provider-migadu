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
	"github.com/metio/terraform-provider-migadu/migadu/client"
)

var (
	_ datasource.DataSource              = &identitiesDataSource{}
	_ datasource.DataSourceWithConfigure = &identitiesDataSource{}
)

func NewIdentitiesDataSource() datasource.DataSource {
	return &identitiesDataSource{}
}

type identitiesDataSource struct {
	migaduClient *client.MigaduClient
}

type identitiesDataSourceModel struct {
	ID         types.String    `tfsdk:"id"`
	DomainName types.String    `tfsdk:"domain_name"`
	LocalPart  types.String    `tfsdk:"local_part"`
	Identities []identityModel `tfsdk:"identities"`
}

type identityModel struct {
	LocalPart            types.String `tfsdk:"local_part"`
	DomainName           types.String `tfsdk:"domain_name"`
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

func (d *identitiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identities"
}

func (d *identitiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gets all identities of a domain.",
		MarkdownDescription: "Gets all identities of a domain.",
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
			"id": schema.StringAttribute{
				Description:         "Contains the full email address 'local_part@domain_name'.",
				MarkdownDescription: "Contains the full email address `local_part@domain_name`.",
				Computed:            true,
			},
			"identities": schema.ListNestedAttribute{
				Description:         "The configured identities for the given 'domain_name' and 'local_part'.",
				MarkdownDescription: "The configured identities for the given `domain_name` and `local_part`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"local_part": schema.StringAttribute{
							Computed: true,
						},
						"domain_name": schema.StringAttribute{
							Computed: true,
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
				},
			},
		},
	}
}

func (d *identitiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *identitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data identitiesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identities, err := d.migaduClient.GetIdentities(ctx, data.DomainName.ValueString(), data.LocalPart.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

	for _, identity := range identities.Identities {
		identityModel := identityModel{
			LocalPart:            types.StringValue(identity.LocalPart),
			DomainName:           types.StringValue(identity.DomainName),
			Address:              types.StringValue(identity.Address),
			Name:                 types.StringValue(identity.Name),
			MaySend:              types.BoolValue(identity.MaySend),
			MayReceive:           types.BoolValue(identity.MayReceive),
			MayAccessImap:        types.BoolValue(identity.MayAccessImap),
			MayAccessPop3:        types.BoolValue(identity.MayAccessPop3),
			MayAccessManageSieve: types.BoolValue(identity.MayAccessManageSieve),
			FooterActive:         types.BoolValue(identity.FooterActive),
			FooterPlainBody:      types.StringValue(identity.FooterPlainBody),
			FooterHtmlBody:       types.StringValue(identity.FooterHtmlBody),
		}

		data.Identities = append(data.Identities, identityModel)
	}

	data.ID = types.StringValue(fmt.Sprintf("%s@%s", data.LocalPart.ValueString(), data.DomainName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
