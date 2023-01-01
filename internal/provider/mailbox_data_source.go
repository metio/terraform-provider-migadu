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
	_ datasource.DataSource              = &mailboxDataSource{}
	_ datasource.DataSourceWithConfigure = &mailboxDataSource{}
)

func NewMailboxDataSource() datasource.DataSource {
	return &mailboxDataSource{}
}

type mailboxDataSource struct {
	migaduClient *client.MigaduClient
}

type mailboxDataSourceModel struct {
	ID                        types.String  `tfsdk:"id"`
	LocalPart                 types.String  `tfsdk:"local_part"`
	DomainName                types.String  `tfsdk:"domain_name"`
	Address                   types.String  `tfsdk:"address"`
	Name                      types.String  `tfsdk:"name"`
	IsInternal                types.Bool    `tfsdk:"is_internal"`
	MaySend                   types.Bool    `tfsdk:"may_send"`
	MayReceive                types.Bool    `tfsdk:"may_receive"`
	MayAccessImap             types.Bool    `tfsdk:"may_access_imap"`
	MayAccessPop3             types.Bool    `tfsdk:"may_access_pop3"`
	MayAccessManageSieve      types.Bool    `tfsdk:"may_access_manage_sieve"`
	PasswordRecoveryEmail     types.String  `tfsdk:"password_recovery_email"`
	SpamAction                types.String  `tfsdk:"spam_action"`
	SpamAggressiveness        types.String  `tfsdk:"spam_aggressiveness"`
	Expirable                 types.Bool    `tfsdk:"expirable"`
	ExpiresOn                 types.String  `tfsdk:"expires_on"`
	RemoveUponExpiry          types.Bool    `tfsdk:"remove_upon_expiry"`
	SenderDenyList            types.List    `tfsdk:"sender_denylist"`
	SenderDenyListPunycode    types.List    `tfsdk:"sender_denylist_punycode"`
	SenderAllowList           types.List    `tfsdk:"sender_allowlist"`
	SenderAllowListPunycode   types.List    `tfsdk:"sender_allowlist_punycode"`
	RecipientDenyList         types.List    `tfsdk:"recipient_denylist"`
	RecipientDenyListPunycode types.List    `tfsdk:"recipient_denylist_punycode"`
	AutoRespondActive         types.Bool    `tfsdk:"auto_respond_active"`
	AutoRespondSubject        types.String  `tfsdk:"auto_respond_subject"`
	AutoRespondBody           types.String  `tfsdk:"auto_respond_body"`
	AutoRespondExpiresOn      types.String  `tfsdk:"auto_respond_expires_on"`
	FooterActive              types.Bool    `tfsdk:"footer_active"`
	FooterPlainBody           types.String  `tfsdk:"footer_plain_body"`
	FooterHtmlBody            types.String  `tfsdk:"footer_html_body"`
	StorageUsage              types.Float64 `tfsdk:"storage_usage"`
	Delegations               types.List    `tfsdk:"delegations"`
	DelegationsPunycode       types.List    `tfsdk:"delegations_punycode"`
	Identities                types.List    `tfsdk:"identities"`
	IdentitiesPunycode        types.List    `tfsdk:"identities_punycode"`
}

func (d *mailboxDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mailbox"
}

func (d *mailboxDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gets a single mailbox.",
		MarkdownDescription: "Gets a single mailbox.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the mailbox to fetch.",
				MarkdownDescription: "The domain name of the mailbox to fetch.",
				Required:            true,
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the mailbox to fetch.",
				MarkdownDescription: "The local part of the mailbox to fetch.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Contains the full email address 'local_part@domain_name'.",
				MarkdownDescription: "Contains the full email address `local_part@domain_name`.",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"is_internal": schema.BoolAttribute{
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
			"password_recovery_email": schema.StringAttribute{
				Computed: true,
			},
			"spam_action": schema.StringAttribute{
				Computed: true,
			},
			"spam_aggressiveness": schema.StringAttribute{
				Computed: true,
			},
			"expirable": schema.BoolAttribute{
				Computed: true,
			},
			"expires_on": schema.StringAttribute{
				Computed: true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Computed: true,
			},
			"sender_denylist": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"sender_denylist_punycode": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"sender_allowlist": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"sender_allowlist_punycode": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"recipient_denylist": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"recipient_denylist_punycode": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"auto_respond_active": schema.BoolAttribute{
				Computed: true,
			},
			"auto_respond_subject": schema.StringAttribute{
				Computed: true,
			},
			"auto_respond_body": schema.StringAttribute{
				Computed: true,
			},
			"auto_respond_expires_on": schema.StringAttribute{
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
			"storage_usage": schema.Float64Attribute{
				Computed: true,
			},
			"delegations": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"delegations_punycode": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"identities": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"identities_punycode": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *mailboxDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *mailboxDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mailboxDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mailbox, err := d.migaduClient.GetMailbox(ctx, data.DomainName.ValueString(), data.LocalPart.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

	senderDenyList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(mailbox.SenderDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderDenyListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(mailbox.SenderDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderAllowList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(mailbox.SenderAllowList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderAllowListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(mailbox.SenderAllowList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	recipientDenyList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(mailbox.RecipientDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	recipientDenyListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(mailbox.RecipientDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	delegations, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(mailbox.Delegations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	delegationsPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(mailbox.Delegations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	identities, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(mailbox.Identities, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	identitiesPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(mailbox.Identities, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s@%s", data.LocalPart.ValueString(), data.DomainName.ValueString()))
	//data.LocalPart = types.StringValue(mailbox.LocalPart)
	//data.DomainName = types.StringValue(mailbox.DomainName)
	data.Address = types.StringValue(mailbox.Address)
	data.Name = types.StringValue(mailbox.Name)
	data.IsInternal = types.BoolValue(mailbox.IsInternal)
	data.MaySend = types.BoolValue(mailbox.MaySend)
	data.MayReceive = types.BoolValue(mailbox.MayReceive)
	data.MayAccessImap = types.BoolValue(mailbox.MayAccessImap)
	data.MayAccessPop3 = types.BoolValue(mailbox.MayAccessPop3)
	data.MayAccessManageSieve = types.BoolValue(mailbox.MayAccessManageSieve)
	data.PasswordRecoveryEmail = types.StringValue(mailbox.PasswordRecoveryEmail)
	data.SpamAction = types.StringValue(mailbox.SpamAction)
	data.SpamAggressiveness = types.StringValue(mailbox.SpamAggressiveness)
	data.Expirable = types.BoolValue(mailbox.Expirable)
	data.ExpiresOn = types.StringValue(mailbox.ExpiresOn)
	data.RemoveUponExpiry = types.BoolValue(mailbox.RemoveUponExpiry)
	data.AutoRespondActive = types.BoolValue(mailbox.AutoRespondActive)
	data.AutoRespondSubject = types.StringValue(mailbox.AutoRespondSubject)
	data.AutoRespondBody = types.StringValue(mailbox.AutoRespondBody)
	data.AutoRespondExpiresOn = types.StringValue(mailbox.AutoRespondExpiresOn)
	data.FooterActive = types.BoolValue(mailbox.FooterActive)
	data.FooterPlainBody = types.StringValue(mailbox.FooterPlainBody)
	data.FooterHtmlBody = types.StringValue(mailbox.FooterHtmlBody)
	data.StorageUsage = types.Float64Value(mailbox.StorageUsage)
	data.SenderDenyList = senderDenyList
	data.SenderDenyListPunycode = senderDenyListPunycode
	data.SenderAllowList = senderAllowList
	data.SenderAllowListPunycode = senderAllowListPunycode
	data.RecipientDenyList = recipientDenyList
	data.RecipientDenyListPunycode = recipientDenyListPunycode
	data.Delegations = delegations
	data.DelegationsPunycode = delegationsPunycode
	data.Identities = identities
	data.IdentitiesPunycode = identitiesPunycode

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
