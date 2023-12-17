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
	_ datasource.DataSource              = (*IdentitiesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*IdentitiesDataSource)(nil)
)

func NewIdentitiesDataSource() datasource.DataSource {
	return &IdentitiesDataSource{}
}

type IdentitiesDataSource struct {
	MigaduClient *client.MigaduClient
}

type IdentitiesDataSourceModel struct {
	ID         custom_types.EmailAddressValue `tfsdk:"id"`
	LocalPart  types.String                   `tfsdk:"local_part"`
	DomainName custom_types.DomainNameValue   `tfsdk:"domain_name"`
	Identities []IdentityModel                `tfsdk:"identities"`
}

type IdentityModel struct {
	LocalPart            types.String                   `tfsdk:"local_part"`
	DomainName           custom_types.DomainNameValue   `tfsdk:"domain_name"`
	Address              custom_types.EmailAddressValue `tfsdk:"address"`
	Name                 types.String                   `tfsdk:"name"`
	MaySend              types.Bool                     `tfsdk:"may_send"`
	MayReceive           types.Bool                     `tfsdk:"may_receive"`
	MayAccessImap        types.Bool                     `tfsdk:"may_access_imap"`
	MayAccessPop3        types.Bool                     `tfsdk:"may_access_pop3"`
	MayAccessManageSieve types.Bool                     `tfsdk:"may_access_manage_sieve"`
	FooterActive         types.Bool                     `tfsdk:"footer_active"`
	FooterPlainBody      types.String                   `tfsdk:"footer_plain_body"`
	FooterHtmlBody       types.String                   `tfsdk:"footer_html_body"`
}

func (d *IdentitiesDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_identities"
}

func (d *IdentitiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Get information about all identities owned by mailbox.",
		MarkdownDescription: "Get information about all identities owned by mailbox.",
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
				Description:         "The local part of the mailbox that owns the identities.",
				MarkdownDescription: "The local part of the mailbox that owns the identities.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the mailbox/identities.",
				MarkdownDescription: "The domain name of the mailbox/identities.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType:          custom_types.DomainNameType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"identities": schema.ListNestedAttribute{
				Description:         "The configured identities for the given 'domain_name' and 'local_part'.",
				MarkdownDescription: "The configured identities for the given `domain_name` and `local_part`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"local_part": schema.StringAttribute{
							Description:         "The local part of the identity.",
							MarkdownDescription: "The local part of the identity.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"domain_name": schema.StringAttribute{
							Description:         "The domain of the identity.",
							MarkdownDescription: "The domain of the identity.",
							Required:            false,
							Optional:            false,
							Computed:            true,
							CustomType:          custom_types.DomainNameType{},
						},
						"address": schema.StringAttribute{
							Description:         "The email address of the identity 'identity@domain_name' as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
							MarkdownDescription: "The email address of the identity `identity@domain_name` as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
							Required:            false,
							Optional:            false,
							Computed:            true,
							CustomType:          custom_types.EmailAddressType{},
						},
						"name": schema.StringAttribute{
							Description:         "The name of the identity.",
							MarkdownDescription: "The name of the identity.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_send": schema.BoolAttribute{
							Description:         "Whether the identity is allowed to send emails.",
							MarkdownDescription: "Whether the identity is allowed to send emails.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_receive": schema.BoolAttribute{
							Description:         "Whether the identity is allowed to receive emails.",
							MarkdownDescription: "Whether the identity is allowed to receive emails.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_access_imap": schema.BoolAttribute{
							Description:         "Whether the identity is allowed to use IMAP.",
							MarkdownDescription: "Whether the identity is allowed to use IMAP.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_access_pop3": schema.BoolAttribute{
							Description:         "Whether the identity is allowed to use POP3.",
							MarkdownDescription: "Whether the identity is allowed to use POP3.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_access_manage_sieve": schema.BoolAttribute{
							Description:         "Whether the identity is allowed to manage the mail sieve.",
							MarkdownDescription: "Whether the identity is allowed to manage the mail sieve.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"footer_active": schema.BoolAttribute{
							Description:         "Whether the footer of the identity is active.",
							MarkdownDescription: "Whether the footer of the identity is active.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"footer_plain_body": schema.StringAttribute{
							Description:         "The footer of the identity in 'text/plain' format.",
							MarkdownDescription: "The footer of the identity in `text/plain` format.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"footer_html_body": schema.StringAttribute{
							Description:         "The footer of the identity in 'text/html' format.",
							MarkdownDescription: "The footer of the identity in `text/html` format.",
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

func (d *IdentitiesDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (d *IdentitiesDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data IdentitiesDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	identities, err := d.MigaduClient.GetIdentities(ctx, data.DomainName.ValueString(), data.LocalPart.ValueString())
	if err != nil {
		response.Diagnostics.Append(IdentityReadError(err))
		return
	}

	for _, identity := range identities.Identities {
		model := IdentityModel{
			LocalPart:            types.StringValue(identity.LocalPart),
			DomainName:           custom_types.NewDomainNameValue(identity.DomainName),
			Address:              custom_types.NewEmailAddressValue(identity.Address),
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
		data.Identities = append(data.Identities, model)
	}

	data.ID = custom_types.NewEmailAddressValue(fmt.Sprintf("%s@%s", data.LocalPart.ValueString(), data.DomainName.ValueString()))

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
