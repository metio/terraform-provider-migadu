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
	ID                        types.String `tfsdk:"id"`
	LocalPart                 types.String `tfsdk:"local_part"`
	DomainName                types.String `tfsdk:"domain_name"`
	Address                   types.String `tfsdk:"address"`
	Name                      types.String `tfsdk:"name"`
	IsInternal                types.Bool   `tfsdk:"is_internal"`
	MaySend                   types.Bool   `tfsdk:"may_send"`
	MayReceive                types.Bool   `tfsdk:"may_receive"`
	MayAccessImap             types.Bool   `tfsdk:"may_access_imap"`
	MayAccessPop3             types.Bool   `tfsdk:"may_access_pop3"`
	MayAccessManageSieve      types.Bool   `tfsdk:"may_access_manage_sieve"`
	PasswordRecoveryEmail     types.String `tfsdk:"password_recovery_email"`
	SpamAction                types.String `tfsdk:"spam_action"`
	SpamAggressiveness        types.String `tfsdk:"spam_aggressiveness"`
	Expirable                 types.Bool   `tfsdk:"expirable"`
	ExpiresOn                 types.String `tfsdk:"expires_on"`
	RemoveUponExpiry          types.Bool   `tfsdk:"remove_upon_expiry"`
	SenderDenyList            types.List   `tfsdk:"sender_denylist"`
	SenderDenyListPunycode    types.List   `tfsdk:"sender_denylist_punycode"`
	SenderAllowList           types.List   `tfsdk:"sender_allowlist"`
	SenderAllowListPunycode   types.List   `tfsdk:"sender_allowlist_punycode"`
	RecipientDenyList         types.List   `tfsdk:"recipient_denylist"`
	RecipientDenyListPunycode types.List   `tfsdk:"recipient_denylist_punycode"`
	AutoRespondActive         types.Bool   `tfsdk:"auto_respond_active"`
	AutoRespondSubject        types.String `tfsdk:"auto_respond_subject"`
	AutoRespondBody           types.String `tfsdk:"auto_respond_body"`
	AutoRespondExpiresOn      types.String `tfsdk:"auto_respond_expires_on"`
	FooterActive              types.Bool   `tfsdk:"footer_active"`
	FooterPlainBody           types.String `tfsdk:"footer_plain_body"`
	FooterHtmlBody            types.String `tfsdk:"footer_html_body"`
	Delegations               types.List   `tfsdk:"delegations"`
	DelegationsPunycode       types.List   `tfsdk:"delegations_punycode"`
	Identities                types.List   `tfsdk:"identities"`
	IdentitiesPunycode        types.List   `tfsdk:"identities_punycode"`
}

func (d *mailboxDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mailbox"
}

func (d *mailboxDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Get information about a mailbox.",
		MarkdownDescription: "Get information about a mailbox.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the mailbox.",
				MarkdownDescription: "The domain name of the mailbox.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the mailbox.",
				MarkdownDescription: "The local part of the mailbox.",
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
				Description:         "The email address of the mailbox 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "The email address of the mailbox `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "The name of the mailbox.",
				MarkdownDescription: "The name of the mailbox.",
				Computed:            true,
			},
			"is_internal": schema.BoolAttribute{
				Description:         "Whether the mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
				MarkdownDescription: "Whether the mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
				Computed:            true,
			},
			"may_send": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to send emails.",
				MarkdownDescription: "Whether the mailbox is allowed to send emails.",
				Computed:            true,
			},
			"may_receive": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to receive emails.",
				MarkdownDescription: "Whether the mailbox is allowed to receive emails.",
				Computed:            true,
			},
			"may_access_imap": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to use IMAP.",
				MarkdownDescription: "Whether the mailbox is allowed to use IMAP.",
				Computed:            true,
			},
			"may_access_pop3": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to use POP3.",
				MarkdownDescription: "Whether the mailbox is allowed to use POP3.",
				Computed:            true,
			},
			"may_access_manage_sieve": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to manage the mail sieve.",
				MarkdownDescription: "Whether the mailbox is allowed to manage the mail sieve.",
				Computed:            true,
			},
			"password_recovery_email": schema.StringAttribute{
				Description:         "The recovery email address of the mailbox.",
				MarkdownDescription: "The recovery email address of the mailbox.",
				Computed:            true,
			},
			"spam_action": schema.StringAttribute{
				Description:         "The action to take once spam arrives in the mailbox.",
				MarkdownDescription: "The action to take once spam arrives in the mailbox.",
				Computed:            true,
			},
			"spam_aggressiveness": schema.StringAttribute{
				Description:         "How aggressive will spam be detected in the mailbox.",
				MarkdownDescription: "How aggressive will spam be detected in the mailbox.",
				Computed:            true,
			},
			"expirable": schema.BoolAttribute{
				Description:         "Whether the mailbox expires in the future.",
				MarkdownDescription: "Whether the mailbox expires in the future.",
				Computed:            true,
			},
			"expires_on": schema.StringAttribute{
				Description:         "The expiration date of the mailbox.",
				MarkdownDescription: "The expiration date of the mailbox.",
				Computed:            true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Description:         "Whether the mailbox will be removed upon expiry.",
				MarkdownDescription: "Whether the mailbox will be removed upon expiry.",
				Computed:            true,
			},
			"sender_denylist": schema.ListAttribute{
				Description:         "The email addresses of senders that will always be denied delivery in unicode.",
				MarkdownDescription: "The email addresses of senders that will always be denied delivery in unicode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sender_denylist_punycode": schema.ListAttribute{
				Description:         "The email addresses of senders that will always be denied delivery in punycode.",
				MarkdownDescription: "The email addresses of senders that will always be denied delivery in punycode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sender_allowlist": schema.ListAttribute{
				Description:         "The email addresses of senders that will always be allowed delivery in unicode.",
				MarkdownDescription: "The email addresses of senders that will always be allowed delivery in unicode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sender_allowlist_punycode": schema.ListAttribute{
				Description:         "The email addresses of senders that will always be denied delivery in punycode.",
				MarkdownDescription: "The email addresses of senders that will always be denied delivery in punycode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"recipient_denylist": schema.ListAttribute{
				Description:         "The email addresses of recipients that will always be denied delivery in unicode.",
				MarkdownDescription: "The email addresses of recipients that will always be denied delivery in unicode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"recipient_denylist_punycode": schema.ListAttribute{
				Description:         "The email addresses of recipients that will always be denied delivery in punycode.",
				MarkdownDescription: "The email addresses of recipients that will always be denied delivery in punycode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"auto_respond_active": schema.BoolAttribute{
				Description:         "Whether an automatic response is active in the mailbox.",
				MarkdownDescription: "Whether an automatic response is active in the mailbox.",
				Computed:            true,
			},
			"auto_respond_subject": schema.StringAttribute{
				Description:         "The subject of the automatic response.",
				MarkdownDescription: "The subject of the automatic response.",
				Computed:            true,
			},
			"auto_respond_body": schema.StringAttribute{
				Description:         "The body of the automatic response.",
				MarkdownDescription: "The body of the automatic response.",
				Computed:            true,
			},
			"auto_respond_expires_on": schema.StringAttribute{
				Description:         "The expiration date of the automatic response.",
				MarkdownDescription: "The expiration date of the automatic response.",
				Computed:            true,
			},
			"footer_active": schema.BoolAttribute{
				Description:         "Whether the footer of the mailbox is active.",
				MarkdownDescription: "Whether the footer of the mailbox is active.",
				Computed:            true,
			},
			"footer_plain_body": schema.StringAttribute{
				Description:         "The footer of the mailbox in 'text/plain' format.",
				MarkdownDescription: "The footer of the mailbox in `text/plain` format.",
				Computed:            true,
			},
			"footer_html_body": schema.StringAttribute{
				Description:         "The footer of the mailbox in 'text/plain' format.",
				MarkdownDescription: "The footer of the mailbox in `text/html` format.",
				Computed:            true,
			},
			"delegations": schema.ListAttribute{
				Description:         "The delegations of the mailbox in unicode.",
				MarkdownDescription: "The delegations of the mailbox in unicode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"delegations_punycode": schema.ListAttribute{
				Description:         "The delegations of the mailbox in punycode.",
				MarkdownDescription: "The delegations of the mailbox in punycode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"identities": schema.ListAttribute{
				Description:         "The identities of the mailbox in unicode.",
				MarkdownDescription: "The identities of the mailbox in unicode.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"identities_punycode": schema.ListAttribute{
				Description:         "The identities of the mailbox in punycode.",
				MarkdownDescription: "The identities of the mailbox in punycode.",
				Computed:            true,
				ElementType:         types.StringType,
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
