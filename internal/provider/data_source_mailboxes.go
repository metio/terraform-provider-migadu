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
	"github.com/metio/terraform-provider-migadu/internal/client"
)

var (
	_ datasource.DataSource              = &MailboxesDataSource{}
	_ datasource.DataSourceWithConfigure = &MailboxesDataSource{}
)

func NewMailboxesDataSource() datasource.DataSource {
	return &MailboxesDataSource{}
}

type MailboxesDataSource struct {
	migaduClient *client.MigaduClient
}

type MailboxesDataSourceModel struct {
	Id         types.String   `tfsdk:"id"`
	DomainName types.String   `tfsdk:"domain_name"`
	Mailboxes  []MailboxModel `tfsdk:"mailboxes"`
}

type MailboxModel struct {
	LocalPart             types.String  `tfsdk:"local_part"`
	DomainName            types.String  `tfsdk:"domain_name"`
	Address               types.String  `tfsdk:"address"`
	Name                  types.String  `tfsdk:"name"`
	IsInternal            types.Bool    `tfsdk:"is_internal"`
	MaySend               types.Bool    `tfsdk:"may_send"`
	MayReceive            types.Bool    `tfsdk:"may_receive"`
	MayAccessImap         types.Bool    `tfsdk:"may_access_imap"`
	MayAccessPop3         types.Bool    `tfsdk:"may_access_pop3"`
	MayAccessManageSieve  types.Bool    `tfsdk:"may_access_managesieve"`
	PasswordRecoveryEmail types.String  `tfsdk:"password_recovery_email"`
	SpamAction            types.String  `tfsdk:"spam_action"`
	SpamAggressiveness    types.String  `tfsdk:"spam_aggressiveness"`
	Expirable             types.Bool    `tfsdk:"expirable"`
	ExpiresOn             types.String  `tfsdk:"expires_on"`
	RemoveUponExpiry      types.Bool    `tfsdk:"remove_upon_expiry"`
	SenderDenyList        types.List    `tfsdk:"sender_denylist"`
	SenderAllowList       types.List    `tfsdk:"sender_allowlist"`
	RecipientDenyList     types.List    `tfsdk:"recipient_denylist"`
	AutoRespondActive     types.Bool    `tfsdk:"autorespond_active"`
	AutoRespondSubject    types.String  `tfsdk:"autorespond_subject"`
	AutoRespondBody       types.String  `tfsdk:"autorespond_body"`
	AutoRespondExpiresOn  types.String  `tfsdk:"autorespond_expires_on"`
	FooterActive          types.Bool    `tfsdk:"footer_active"`
	FooterPlainBody       types.String  `tfsdk:"footer_plain_body"`
	FooterHtmlBody        types.String  `tfsdk:"footer_html_body"`
	StorageUsage          types.Float64 `tfsdk:"storage_usage"`
	Delegations           types.List    `tfsdk:"delegations"`
	Identities            types.List    `tfsdk:"identities"`
}

func (d *MailboxesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mailboxes"
}

func (d *MailboxesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gets all mailboxes of a domain.",
		MarkdownDescription: "Gets all mailboxes of a domain.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain to fetch mailboxes of.",
				MarkdownDescription: "The domain to fetch mailboxes of.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Same value as the 'domain_name' attribute.",
				MarkdownDescription: "Same value as the `domain_name` attribute.",
				Computed:            true,
			},
			"mailboxes": schema.ListNestedAttribute{
				Description:         "The configured mailboxes for the given 'domain_name'.",
				MarkdownDescription: "The configured mailboxes for the given `domain_name`.",
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
						"may_access_managesieve": schema.BoolAttribute{
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
						"sender_allowlist": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"recipient_denylist": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"autorespond_active": schema.BoolAttribute{
							Computed: true,
						},
						"autorespond_subject": schema.StringAttribute{
							Computed: true,
						},
						"autorespond_body": schema.StringAttribute{
							Computed: true,
						},
						"autorespond_expires_on": schema.StringAttribute{
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
						"identities": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *MailboxesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MailboxesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MailboxesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mailboxes, err := d.migaduClient.GetMailboxes(data.DomainName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Migadu Client Error", "Request failed with: "+err.Error())
		return
	}

	for _, mailbox := range mailboxes.Mailboxes {
		mailboxModel := MailboxModel{
			LocalPart:             types.StringValue(mailbox.LocalPart),
			DomainName:            types.StringValue(mailbox.DomainName),
			Address:               types.StringValue(mailbox.Address),
			Name:                  types.StringValue(mailbox.Name),
			IsInternal:            types.BoolValue(mailbox.IsInternal),
			MaySend:               types.BoolValue(mailbox.MaySend),
			MayReceive:            types.BoolValue(mailbox.MayReceive),
			MayAccessImap:         types.BoolValue(mailbox.MayAccessImap),
			MayAccessPop3:         types.BoolValue(mailbox.MayAccessPop3),
			MayAccessManageSieve:  types.BoolValue(mailbox.MayAccessManageSieve),
			PasswordRecoveryEmail: types.StringValue(mailbox.PasswordRecoveryEmail),
			SpamAction:            types.StringValue(mailbox.SpamAction),
			SpamAggressiveness:    types.StringValue(mailbox.SpamAggressiveness),
			Expirable:             types.BoolValue(mailbox.Expirable),
			ExpiresOn:             types.StringValue(mailbox.ExpiresOn),
			RemoveUponExpiry:      types.BoolValue(mailbox.RemoveUponExpiry),
			AutoRespondActive:     types.BoolValue(mailbox.AutoRespondActive),
			AutoRespondSubject:    types.StringValue(mailbox.AutoRespondSubject),
			AutoRespondBody:       types.StringValue(mailbox.AutoRespondBody),
			AutoRespondExpiresOn:  types.StringValue(mailbox.AutoRespondExpiresOn),
			FooterActive:          types.BoolValue(mailbox.FooterActive),
			FooterPlainBody:       types.StringValue(mailbox.FooterPlainBody),
			FooterHtmlBody:        types.StringValue(mailbox.FooterHtmlBody),
			StorageUsage:          types.Float64Value(mailbox.StorageUsage),
		}

		senderDenyList, diags := types.ListValueFrom(ctx, types.StringType, mailbox.SenderDenyList)
		resp.Diagnostics.Append(diags...)
		senderAllowList, diags := types.ListValueFrom(ctx, types.StringType, mailbox.SenderAllowList)
		resp.Diagnostics.Append(diags...)
		recipientDenyList, diags := types.ListValueFrom(ctx, types.StringType, mailbox.RecipientDenyList)
		resp.Diagnostics.Append(diags...)
		delegations, diags := types.ListValueFrom(ctx, types.StringType, mailbox.Delegations)
		resp.Diagnostics.Append(diags...)
		identities, diags := types.ListValueFrom(ctx, types.StringType, mailbox.Identities)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		mailboxModel.SenderDenyList = senderDenyList
		mailboxModel.SenderAllowList = senderAllowList
		mailboxModel.RecipientDenyList = recipientDenyList
		mailboxModel.Delegations = delegations
		mailboxModel.Identities = identities

		data.Mailboxes = append(data.Mailboxes, mailboxModel)
	}

	data.Id = data.DomainName

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
