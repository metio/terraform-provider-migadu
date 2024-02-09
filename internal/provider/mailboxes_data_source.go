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
	_ datasource.DataSource              = (*MailboxesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*MailboxesDataSource)(nil)
)

func NewMailboxesDataSource() datasource.DataSource {
	return &MailboxesDataSource{}
}

type MailboxesDataSource struct {
	migaduClient *client.MigaduClient
}

type MailboxesDataSourceModel struct {
	ID         custom_types.DomainNameValue `tfsdk:"id"`
	DomainName custom_types.DomainNameValue `tfsdk:"domain_name"`
	Mailboxes  []MailboxModel               `tfsdk:"mailboxes"`
}

type MailboxModel struct {
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

func (d *MailboxesDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mailboxes"
}

func (d *MailboxesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Get information about all mailbox of a domain.",
		MarkdownDescription: "Get information about all mailbox of a domain.",
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
				Description:         "The domain name of the mailboxes.",
				MarkdownDescription: "The domain name of the mailboxes.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType:          custom_types.DomainNameType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"mailboxes": schema.ListNestedAttribute{
				Description:         "The configured mailboxes for the given 'domain_name'.",
				MarkdownDescription: "The configured mailboxes for the given `domain_name`.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"local_part": schema.StringAttribute{
							Description:         "The local part of the mailbox.",
							MarkdownDescription: "The local part of the mailbox.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"domain_name": schema.StringAttribute{
							Description:         "The domain name of the mailbox.",
							MarkdownDescription: "The domain name of the mailbox.",
							Required:            false,
							Optional:            false,
							Computed:            true,
							CustomType:          custom_types.DomainNameType{},
						},
						"address": schema.StringAttribute{
							Description:         "The email address of the mailbox 'local_part@domain_name' as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
							MarkdownDescription: "The email address of the mailbox `local_part@domain_name` as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
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
							Description:         "Whether this mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
							MarkdownDescription: "Whether this mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_send": schema.BoolAttribute{
							Description:         "Whether this mailbox is allowed to send emails.",
							MarkdownDescription: "Whether this mailbox is allowed to send emails.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_receive": schema.BoolAttribute{
							Description:         "Whether this mailbox is allowed to receive emails.",
							MarkdownDescription: "Whether this mailbox is allowed to receive emails.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_access_imap": schema.BoolAttribute{
							Description:         "Whether this mailbox is allowed to use IMAP.",
							MarkdownDescription: "Whether this mailbox is allowed to use IMAP.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_access_pop3": schema.BoolAttribute{
							Description:         "Whether this mailbox is allowed to use POP3.",
							MarkdownDescription: "Whether this mailbox is allowed to use POP3.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"may_access_manage_sieve": schema.BoolAttribute{
							Description:         "Whether this mailbox is allowed to manage the mail sieve.",
							MarkdownDescription: "Whether this mailbox is allowed to manage the mail sieve.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"password_recovery_email": schema.StringAttribute{
							Description:         "The recovery email address of this mailbox.",
							MarkdownDescription: "The recovery email address of this mailbox.",
							Required:            false,
							Optional:            false,
							Computed:            true,
							CustomType:          custom_types.EmailAddressType{},
						},
						"spam_action": schema.StringAttribute{
							Description:         "The action to take once spam arrives in this mailbox.",
							MarkdownDescription: "The action to take once spam arrives in this mailbox.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"spam_aggressiveness": schema.StringAttribute{
							Description:         "How aggressive will spam be detected in this mailbox.",
							MarkdownDescription: "How aggressive will spam be detected in this mailbox.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"expirable": schema.BoolAttribute{
							Description:         "Whether this mailbox expires in the future.",
							MarkdownDescription: "Whether this mailbox expires in the future.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"expires_on": schema.StringAttribute{
							Description:         "The expiration date of this mailbox.",
							MarkdownDescription: "The expiration date of this mailbox.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"remove_upon_expiry": schema.BoolAttribute{
							Description:         "Whether this mailbox will be removed upon expiry.",
							MarkdownDescription: "Whether this mailbox will be removed upon expiry.",
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
							Description:         "Whether an automatic response is active in this mailbox.",
							MarkdownDescription: "Whether an automatic response is active in this mailbox.",
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
							Description:         "Whether the footer of this mailbox is active.",
							MarkdownDescription: "Whether the footer of this mailbox is active.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"footer_plain_body": schema.StringAttribute{
							Description:         "The footer of this mailbox in text/plain format.",
							MarkdownDescription: "The footer of this mailbox in text/plain format.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"footer_html_body": schema.StringAttribute{
							Description:         "The footer of this mailbox in text/html format.",
							MarkdownDescription: "The footer of this mailbox in text/html format.",
							Required:            false,
							Optional:            false,
							Computed:            true,
						},
						"delegations": schema.SetAttribute{
							Description:         "The delegations of this mailbox.",
							MarkdownDescription: "The delegations of this mailbox.",
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
				},
			},
		},
	}
}

func (d *MailboxesDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if migaduClient, ok := request.ProviderData.(*client.MigaduClient); ok {
		d.migaduClient = migaduClient
	} else {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.MigaduClient, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
}

func (d *MailboxesDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data MailboxesDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	mailboxes, err := d.migaduClient.GetMailboxes(ctx, data.DomainName.ValueString())
	if err != nil {
		response.Diagnostics.Append(MailboxReadError(err))
		return
	}

	for _, mailbox := range mailboxes.Mailboxes {
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

		model := MailboxModel{
			LocalPart:             types.StringValue(mailbox.LocalPart),
			DomainName:            custom_types.NewDomainNameValue(mailbox.DomainName),
			Address:               custom_types.NewEmailAddressValue(mailbox.Address),
			Name:                  types.StringValue(mailbox.Name),
			IsInternal:            types.BoolValue(mailbox.IsInternal),
			MaySend:               types.BoolValue(mailbox.MaySend),
			MayReceive:            types.BoolValue(mailbox.MayReceive),
			MayAccessImap:         types.BoolValue(mailbox.MayAccessImap),
			MayAccessPop3:         types.BoolValue(mailbox.MayAccessPop3),
			MayAccessManageSieve:  types.BoolValue(mailbox.MayAccessManageSieve),
			PasswordRecoveryEmail: custom_types.NewEmailAddressValue(mailbox.PasswordRecoveryEmail),
			SpamAction:            types.StringValue(mailbox.SpamAction),
			SpamAggressiveness:    types.StringValue(mailbox.SpamAggressiveness),
			Expirable:             types.BoolValue(mailbox.Expirable),
			ExpiresOn:             types.StringValue(mailbox.ExpiresOn),
			RemoveUponExpiry:      types.BoolValue(mailbox.RemoveUponExpiry),
			SenderDenyList:        senderDenyList,
			SenderAllowList:       senderAllowList,
			RecipientDenyList:     recipientDenyList,
			AutoRespondActive:     types.BoolValue(mailbox.AutoRespondActive),
			AutoRespondSubject:    types.StringValue(mailbox.AutoRespondSubject),
			AutoRespondBody:       types.StringValue(mailbox.AutoRespondBody),
			AutoRespondExpiresOn:  types.StringValue(mailbox.AutoRespondExpiresOn),
			FooterActive:          types.BoolValue(mailbox.FooterActive),
			FooterPlainBody:       types.StringValue(mailbox.FooterPlainBody),
			FooterHtmlBody:        types.StringValue(mailbox.FooterHtmlBody),
			Delegations:           delegations,
		}
		data.Mailboxes = append(data.Mailboxes, model)
	}

	data.ID = data.DomainName

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
