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
	_ datasource.DataSource              = (*MailboxDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*MailboxDataSource)(nil)
)

func NewMailboxDataSource() datasource.DataSource {
	return &MailboxDataSource{}
}

type MailboxDataSource struct {
	MigaduClient *client.MigaduClient
}

type MailboxDataSourceModel struct {
	ID                    custom_types.EmailAddressValue    `tfsdk:"id"`
	LocalPart             types.String                      `tfsdk:"local_part"`
	DomainName            custom_types.DomainNameValue      `tfsdk:"domain_name"`
	Address               custom_types.EmailAddressValue    `tfsdk:"address"`
	Name                  types.String                      `tfsdk:"name"`
	IsInternal            types.Bool                        `tfsdk:"is_internal"`
	MaySend               types.Bool                        `tfsdk:"may_send"`
	MayReceive            types.Bool                        `tfsdk:"may_receive"`
	MayAccessImap         types.Bool                        `tfsdk:"may_access_imap"`
	MayAccessPop3         types.Bool                        `tfsdk:"may_access_pop3"`
	MayAccessManageSieve  types.Bool                        `tfsdk:"may_access_manage_sieve"`
	PasswordRecoveryEmail custom_types.EmailAddressValue    `tfsdk:"password_recovery_email"`
	SpamAction            types.String                      `tfsdk:"spam_action"`
	SpamAggressiveness    types.String                      `tfsdk:"spam_aggressiveness"`
	Expirable             types.Bool                        `tfsdk:"expirable"`
	ExpiresOn             types.String                      `tfsdk:"expires_on"`
	RemoveUponExpiry      types.Bool                        `tfsdk:"remove_upon_expiry"`
	SenderDenyList        custom_types.EmailAddressSetValue `tfsdk:"sender_denylist"`
	SenderAllowList       custom_types.EmailAddressSetValue `tfsdk:"sender_allowlist"`
	RecipientDenyList     custom_types.EmailAddressSetValue `tfsdk:"recipient_denylist"`
	AutoRespondActive     types.Bool                        `tfsdk:"auto_respond_active"`
	AutoRespondSubject    types.String                      `tfsdk:"auto_respond_subject"`
	AutoRespondBody       types.String                      `tfsdk:"auto_respond_body"`
	AutoRespondExpiresOn  types.String                      `tfsdk:"auto_respond_expires_on"`
	FooterActive          types.Bool                        `tfsdk:"footer_active"`
	FooterPlainBody       types.String                      `tfsdk:"footer_plain_body"`
	FooterHtmlBody        types.String                      `tfsdk:"footer_html_body"`
	Delegations           custom_types.EmailAddressSetValue `tfsdk:"delegations"`
}

func (d *MailboxDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mailbox"
}

func (d *MailboxDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Get information about a mailbox.",
		MarkdownDescription: "Get information about a mailbox.",
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
				Description:         "The local part of the mailbox.",
				MarkdownDescription: "The local part of the mailbox.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the mailbox.",
				MarkdownDescription: "The domain name of the mailbox.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType:          custom_types.DomainNameType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"address": schema.StringAttribute{
				Description:         "The email address of the mailbox 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "The email address of the mailbox `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the mailbox.",
				MarkdownDescription: "The name of the mailbox.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"is_internal": schema.BoolAttribute{
				Description:         "Whether the mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
				MarkdownDescription: "Whether the mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"may_send": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to send emails.",
				MarkdownDescription: "Whether the mailbox is allowed to send emails.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"may_receive": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to receive emails.",
				MarkdownDescription: "Whether the mailbox is allowed to receive emails.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"may_access_imap": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to use IMAP.",
				MarkdownDescription: "Whether the mailbox is allowed to use IMAP.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"may_access_pop3": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to use POP3.",
				MarkdownDescription: "Whether the mailbox is allowed to use POP3.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"may_access_manage_sieve": schema.BoolAttribute{
				Description:         "Whether the mailbox is allowed to manage the mail sieve.",
				MarkdownDescription: "Whether the mailbox is allowed to manage the mail sieve.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"password_recovery_email": schema.StringAttribute{
				Description:         "The recovery email address of the mailbox.",
				MarkdownDescription: "The recovery email address of the mailbox.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
			},
			"spam_action": schema.StringAttribute{
				Description:         "The action to take once spam arrives in the mailbox.",
				MarkdownDescription: "The action to take once spam arrives in the mailbox.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"spam_aggressiveness": schema.StringAttribute{
				Description:         "How aggressive will spam be detected in the mailbox.",
				MarkdownDescription: "How aggressive will spam be detected in the mailbox.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"expirable": schema.BoolAttribute{
				Description:         "Whether the mailbox expires in the future.",
				MarkdownDescription: "Whether the mailbox expires in the future.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"expires_on": schema.StringAttribute{
				Description:         "The expiration date of the mailbox.",
				MarkdownDescription: "The expiration date of the mailbox.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Description:         "Whether the mailbox will be removed upon expiry.",
				MarkdownDescription: "Whether the mailbox will be removed upon expiry.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"sender_denylist": schema.SetAttribute{
				Description:         "The email addresses of senders that will always be denied delivery.",
				MarkdownDescription: "The email addresses of senders that will always be denied delivery.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType: custom_types.EmailAddressSetType{
					SetType: types.SetType{
						ElemType: custom_types.EmailAddressType{},
					},
				},
			},
			"sender_allowlist": schema.SetAttribute{
				Description:         "The email addresses of senders that will always be allowed delivery.",
				MarkdownDescription: "The email addresses of senders that will always be allowed delivery.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType: custom_types.EmailAddressSetType{
					SetType: types.SetType{
						ElemType: custom_types.EmailAddressType{},
					},
				},
			},
			"recipient_denylist": schema.SetAttribute{
				Description:         "The email addresses of recipients that will always be denied delivery.",
				MarkdownDescription: "The email addresses of recipients that will always be denied delivery.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType: custom_types.EmailAddressSetType{
					SetType: types.SetType{
						ElemType: custom_types.EmailAddressType{},
					},
				},
			},
			"auto_respond_active": schema.BoolAttribute{
				Description:         "Whether an automatic response is active in the mailbox.",
				MarkdownDescription: "Whether an automatic response is active in the mailbox.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"auto_respond_subject": schema.StringAttribute{
				Description:         "The subject of the automatic response.",
				MarkdownDescription: "The subject of the automatic response.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"auto_respond_body": schema.StringAttribute{
				Description:         "The body of the automatic response.",
				MarkdownDescription: "The body of the automatic response.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"auto_respond_expires_on": schema.StringAttribute{
				Description:         "The expiration date of the automatic response.",
				MarkdownDescription: "The expiration date of the automatic response.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"footer_active": schema.BoolAttribute{
				Description:         "Whether the footer of the mailbox is active.",
				MarkdownDescription: "Whether the footer of the mailbox is active.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"footer_plain_body": schema.StringAttribute{
				Description:         "The footer of the mailbox in 'text/plain' format.",
				MarkdownDescription: "The footer of the mailbox in `text/plain` format.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"footer_html_body": schema.StringAttribute{
				Description:         "The footer of the mailbox in 'text/plain' format.",
				MarkdownDescription: "The footer of the mailbox in `text/html` format.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"delegations": schema.SetAttribute{
				Description:         "The delegations of the mailbox.",
				MarkdownDescription: "The delegations of the mailbox.",
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
	}
}

func (d *MailboxDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (d *MailboxDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data MailboxDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	mailbox, err := d.MigaduClient.GetMailbox(ctx, data.DomainName.ValueString(), data.LocalPart.ValueString())
	if err != nil {
		response.Diagnostics.Append(MailboxReadError(err))
		return
	}

	senderDenyList, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.SenderDenyList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	senderAllowList, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.SenderAllowList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	recipientDenyList, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.RecipientDenyList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	delegations, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.Delegations)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data.ID = custom_types.NewEmailAddressValue(fmt.Sprintf("%s@%s", data.LocalPart.ValueString(), data.DomainName.ValueString()))
	data.Address = custom_types.NewEmailAddressValue(mailbox.Address)
	data.Name = types.StringValue(mailbox.Name)
	data.IsInternal = types.BoolValue(mailbox.IsInternal)
	data.MaySend = types.BoolValue(mailbox.MaySend)
	data.MayReceive = types.BoolValue(mailbox.MayReceive)
	data.MayAccessImap = types.BoolValue(mailbox.MayAccessImap)
	data.MayAccessPop3 = types.BoolValue(mailbox.MayAccessPop3)
	data.MayAccessManageSieve = types.BoolValue(mailbox.MayAccessManageSieve)
	data.PasswordRecoveryEmail = custom_types.NewEmailAddressValue(mailbox.PasswordRecoveryEmail)
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
	data.SenderAllowList = senderAllowList
	data.RecipientDenyList = recipientDenyList
	data.Delegations = delegations

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
