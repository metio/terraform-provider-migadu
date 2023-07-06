/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-migadu/migadu/client"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"strings"
)

var (
	_ resource.Resource                = &mailboxResource{}
	_ resource.ResourceWithConfigure   = &mailboxResource{}
	_ resource.ResourceWithImportState = &mailboxResource{}
)

func NewMailboxResource() resource.Resource {
	return &mailboxResource{}
}

type mailboxResource struct {
	migaduClient *client.MigaduClient
}

type mailboxResourceModel struct {
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
	Password                  types.String `tfsdk:"password"`
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
	Delegations               types.List   `tfsdk:"delegations"`
	DelegationsPunycode       types.List   `tfsdk:"delegations_punycode"`
	Identities                types.List   `tfsdk:"identities"`
	IdentitiesPunycode        types.List   `tfsdk:"identities_punycode"`
	AutoRespondActive         types.Bool   `tfsdk:"auto_respond_active"`
	AutoRespondSubject        types.String `tfsdk:"auto_respond_subject"`
	AutoRespondBody           types.String `tfsdk:"auto_respond_body"`
	AutoRespondExpiresOn      types.String `tfsdk:"auto_respond_expires_on"`
	FooterActive              types.Bool   `tfsdk:"footer_active"`
	FooterPlainBody           types.String `tfsdk:"footer_plain_body"`
	FooterHtmlBody            types.String `tfsdk:"footer_html_body"`
}

func (r *mailboxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mailbox"
}

func (r *mailboxResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Provides a mailbox.",
		MarkdownDescription: "Provides a mailbox.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the mailbox.",
				MarkdownDescription: "The domain name of the mailbox.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the mailbox.",
				MarkdownDescription: "The local part of the mailbox.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Contains the value 'local_part@domain_name'.",
				MarkdownDescription: "Contains the value `local_part@domain_name`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				Description:         "The email address of the mailbox 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "The email address of the mailbox `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the mailbox.",
				MarkdownDescription: "The name of the mailbox.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"is_internal": schema.BoolAttribute{
				Description:         "Whether this mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
				MarkdownDescription: "Whether this mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
				Optional:            true,
				Computed:            true,
			},
			"may_send": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to send emails.",
				MarkdownDescription: "Whether this mailbox is allowed to send emails.",
				Optional:            true,
				Computed:            true,
			},
			"may_receive": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to receive emails.",
				MarkdownDescription: "Whether this mailbox is allowed to receive emails.",
				Optional:            true,
				Computed:            true,
			},
			"may_access_imap": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to use IMAP.",
				MarkdownDescription: "Whether this mailbox is allowed to use IMAP.",
				Optional:            true,
				Computed:            true,
			},
			"may_access_pop3": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to use POP3.",
				MarkdownDescription: "Whether this mailbox is allowed to use POP3.",
				Optional:            true,
				Computed:            true,
			},
			"may_access_manage_sieve": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to manage the mail sieve.",
				MarkdownDescription: "Whether this mailbox is allowed to manage the mail sieve.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				Description:         "The password of this mailbox.",
				MarkdownDescription: "The password of this mailbox.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("password_recovery_email")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password_recovery_email": schema.StringAttribute{
				Description:         "The recovery email address of this mailbox. If this is set instead of 'password' an invitation to that address will be send to the user and they can set their own password.",
				MarkdownDescription: "The recovery email address of this mailbox. If this is set instead of `password` an invitation to that address will be send to the user and they can set their own password.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("password")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"spam_action": schema.StringAttribute{
				Description:         "The action to take once spam arrives in this mailbox.",
				MarkdownDescription: "The action to take once spam arrives in this mailbox.",
				Optional:            true,
				Computed:            true,
			},
			"spam_aggressiveness": schema.StringAttribute{
				Description:         "How aggressive will spam be detected in this mailbox.",
				MarkdownDescription: "How aggressive will spam be detected in this mailbox.",
				Optional:            true,
				Computed:            true,
			},
			"expirable": schema.BoolAttribute{
				Description:         "Whether this mailbox expires in the future.",
				MarkdownDescription: "Whether this mailbox expires in the future.",
				Optional:            true,
				Computed:            true,
			},
			"expires_on": schema.StringAttribute{
				Description:         "The expiration date of this mailbox.",
				MarkdownDescription: "The expiration date of this mailbox.",
				Optional:            true,
				Computed:            true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Description:         "Whether this mailbox will be removed upon expiry.",
				MarkdownDescription: "Whether this mailbox will be removed upon expiry.",
				Optional:            true,
				Computed:            true,
			},
			"sender_denylist": schema.ListAttribute{
				Description:         "The email addresses of senders that will always be denied delivery in unicode.",
				MarkdownDescription: "The email addresses of senders that will always be denied delivery in unicode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sender_denylist_punycode": schema.ListAttribute{
				Description:         "The email addresses of senders that will always be denied delivery in punycode.",
				MarkdownDescription: "The email addresses of senders that will always be denied delivery in punycode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sender_allowlist": schema.ListAttribute{
				Description:         "The email addresses of senders that will always be allowed delivery in unicode.",
				MarkdownDescription: "The email addresses of senders that will always be allowed delivery in unicode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sender_allowlist_punycode": schema.ListAttribute{
				Description:         "The email addresses of senders that will always be denied delivery in punycode.",
				MarkdownDescription: "The email addresses of senders that will always be denied delivery in punycode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"recipient_denylist": schema.ListAttribute{
				Description:         "The email addresses of recipients that will always be denied delivery in unicode.",
				MarkdownDescription: "The email addresses of recipients that will always be denied delivery in unicode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"recipient_denylist_punycode": schema.ListAttribute{
				Description:         "The email addresses of recipients that will always be denied delivery in punycode.",
				MarkdownDescription: "The email addresses of recipients that will always be denied delivery in punycode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"delegations": schema.ListAttribute{
				Description:         "The delegations of the mailbox in unicode.",
				MarkdownDescription: "The delegations of the mailbox in unicode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"delegations_punycode": schema.ListAttribute{
				Description:         "The delegations of the mailbox in punycode.",
				MarkdownDescription: "The delegations of the mailbox in punycode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"identities": schema.ListAttribute{
				Description:         "The identities of the mailbox in unicode.",
				MarkdownDescription: "The identities of the mailbox in unicode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"identities_punycode": schema.ListAttribute{
				Description:         "The identities of the mailbox in punycode.",
				MarkdownDescription: "The identities of the mailbox in punycode.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"auto_respond_active": schema.BoolAttribute{
				Description:         "Whether an automatic response is active in this mailbox.",
				MarkdownDescription: "Whether an automatic response is active in this mailbox.",
				Optional:            true,
				Computed:            true,
			},
			"auto_respond_subject": schema.StringAttribute{
				Description:         "The subject of the automatic response.",
				MarkdownDescription: "The subject of the automatic response.",
				Optional:            true,
				Computed:            true,
			},
			"auto_respond_body": schema.StringAttribute{
				Description:         "The body of the automatic response.",
				MarkdownDescription: "The body of the automatic response.",
				Optional:            true,
				Computed:            true,
			},
			"auto_respond_expires_on": schema.StringAttribute{
				Description:         "The expiration date of the automatic response.",
				MarkdownDescription: "The expiration date of the automatic response.",
				Optional:            true,
				Computed:            true,
			},
			"footer_active": schema.BoolAttribute{
				Description:         "Whether the footer of this mailbox is active.",
				MarkdownDescription: "Whether the footer of this mailbox is active.",
				Optional:            true,
				Computed:            true,
			},
			"footer_plain_body": schema.StringAttribute{
				Description:         "The footer of this mailbox in text/plain format.",
				MarkdownDescription: "The footer of this mailbox in text/plain format.",
				Optional:            true,
				Computed:            true,
			},
			"footer_html_body": schema.StringAttribute{
				Description:         "The footer of this mailbox in text/html format.",
				MarkdownDescription: "The footer of this mailbox in text/html format.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}
func (r *mailboxResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	migaduClient, ok := req.ProviderData.(*client.MigaduClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.MigaduClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.migaduClient = migaduClient
}

func (r *mailboxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan mailboxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var wantedSenderDenyList []string
	if !plan.SenderDenyList.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderDenyList.ElementsAs(ctx, &wantedSenderDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.SenderDenyListPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderDenyListPunycode.ElementsAs(ctx, &wantedSenderDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var wantedSenderAllowList []string
	if !plan.SenderAllowList.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderAllowList.ElementsAs(ctx, &wantedSenderAllowList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.SenderAllowListPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderAllowListPunycode.ElementsAs(ctx, &wantedSenderAllowList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var wantedRecipientDenyList []string
	if !plan.RecipientDenyList.IsUnknown() {
		resp.Diagnostics.Append(plan.RecipientDenyList.ElementsAs(ctx, &wantedRecipientDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.RecipientDenyListPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.RecipientDenyListPunycode.ElementsAs(ctx, &wantedRecipientDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var wantedDelegations []string
	if !plan.Delegations.IsUnknown() {
		resp.Diagnostics.Append(plan.Delegations.ElementsAs(ctx, &wantedDelegations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.DelegationsPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.DelegationsPunycode.ElementsAs(ctx, &wantedDelegations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var wantedIdentities []string
	if !plan.Identities.IsUnknown() {
		resp.Diagnostics.Append(plan.Identities.ElementsAs(ctx, &wantedIdentities, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.IdentitiesPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.IdentitiesPunycode.ElementsAs(ctx, &wantedIdentities, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	mailbox := &model.Mailbox{
		LocalPart:             plan.LocalPart.ValueString(),
		Name:                  plan.Name.ValueString(),
		IsInternal:            plan.IsInternal.ValueBool(),
		MaySend:               plan.MaySend.ValueBool(),
		MayReceive:            plan.MayReceive.ValueBool(),
		MayAccessImap:         plan.MayAccessImap.ValueBool(),
		MayAccessPop3:         plan.MayAccessPop3.ValueBool(),
		MayAccessManageSieve:  plan.MayAccessManageSieve.ValueBool(),
		Password:              plan.Password.ValueString(),
		PasswordRecoveryEmail: plan.PasswordRecoveryEmail.ValueString(),
		SpamAction:            plan.SpamAction.ValueString(),
		SpamAggressiveness:    plan.SpamAggressiveness.ValueString(),
		Expirable:             plan.Expirable.ValueBool(),
		ExpiresOn:             plan.ExpiresOn.ValueString(),
		RemoveUponExpiry:      plan.RemoveUponExpiry.ValueBool(),
		SenderDenyList:        wantedSenderDenyList,
		SenderAllowList:       wantedSenderAllowList,
		RecipientDenyList:     wantedRecipientDenyList,
		Delegations:           wantedDelegations,
		Identities:            wantedIdentities,
		AutoRespondActive:     plan.AutoRespondActive.ValueBool(),
		AutoRespondSubject:    plan.AutoRespondSubject.ValueString(),
		AutoRespondBody:       plan.AutoRespondBody.ValueString(),
		AutoRespondExpiresOn:  plan.AutoRespondExpiresOn.ValueString(),
		FooterActive:          plan.FooterActive.ValueBool(),
		FooterPlainBody:       plan.FooterPlainBody.ValueString(),
		FooterHtmlBody:        plan.FooterHtmlBody.ValueString(),
	}

	if plan.Password.ValueString() != "" {
		mailbox.PasswordMethod = "password"
	} else if plan.PasswordRecoveryEmail.ValueString() != "" {
		mailbox.PasswordMethod = "invitation"
	}

	createdMailbox, err := r.migaduClient.CreateMailbox(ctx, plan.DomainName.ValueString(), mailbox)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating mailbox",
			fmt.Sprintf("Could not create mailbox %s: %v", createMailboxID(plan.DomainName, plan.LocalPart), err),
		)
		return
	}

	senderDenyList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(createdMailbox.SenderDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderDenyListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(createdMailbox.SenderDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderAllowList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(createdMailbox.SenderAllowList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderAllowListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(createdMailbox.SenderAllowList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	recipientDenyList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(createdMailbox.RecipientDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	recipientDenyListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(createdMailbox.RecipientDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	delegations, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(createdMailbox.Delegations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	delegationsPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(createdMailbox.Delegations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	identities, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(createdMailbox.Identities, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	identitiesPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(createdMailbox.Identities, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)

	plan.ID = types.StringValue(createMailboxID(plan.DomainName, plan.LocalPart))
	plan.Address = types.StringValue(createdMailbox.Address)
	plan.Name = types.StringValue(createdMailbox.Name)
	plan.LocalPart = types.StringValue(createdMailbox.LocalPart)
	plan.IsInternal = types.BoolValue(createdMailbox.IsInternal)
	plan.MaySend = types.BoolValue(createdMailbox.MaySend)
	plan.MayReceive = types.BoolValue(createdMailbox.MayReceive)
	plan.MayAccessImap = types.BoolValue(createdMailbox.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(createdMailbox.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(createdMailbox.MayAccessManageSieve)
	plan.Password = types.StringValue(plan.Password.ValueString())
	plan.PasswordRecoveryEmail = types.StringValue(createdMailbox.PasswordRecoveryEmail)
	plan.SpamAction = types.StringValue(createdMailbox.SpamAction)
	plan.SpamAggressiveness = types.StringValue(createdMailbox.SpamAggressiveness)
	plan.Expirable = types.BoolValue(createdMailbox.Expirable)
	plan.ExpiresOn = types.StringValue(createdMailbox.ExpiresOn)
	plan.RemoveUponExpiry = types.BoolValue(createdMailbox.RemoveUponExpiry)
	plan.SenderDenyList = senderDenyList
	plan.SenderDenyListPunycode = senderDenyListPunycode
	plan.SenderAllowList = senderAllowList
	plan.SenderAllowListPunycode = senderAllowListPunycode
	plan.RecipientDenyList = recipientDenyList
	plan.RecipientDenyListPunycode = recipientDenyListPunycode
	plan.Delegations = delegations
	plan.DelegationsPunycode = delegationsPunycode
	plan.Identities = identities
	plan.IdentitiesPunycode = identitiesPunycode
	plan.AutoRespondActive = types.BoolValue(createdMailbox.AutoRespondActive)
	plan.AutoRespondSubject = types.StringValue(createdMailbox.AutoRespondSubject)
	plan.AutoRespondBody = types.StringValue(createdMailbox.AutoRespondBody)
	plan.AutoRespondExpiresOn = types.StringValue(createdMailbox.AutoRespondExpiresOn)
	plan.FooterActive = types.BoolValue(createdMailbox.FooterActive)
	plan.FooterPlainBody = types.StringValue(createdMailbox.FooterPlainBody)
	plan.FooterHtmlBody = types.StringValue(createdMailbox.FooterHtmlBody)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *mailboxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state mailboxResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mailbox, err := r.migaduClient.GetMailbox(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString())
	if err != nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Could not read mailbox %s", createMailboxID(state.DomainName, state.LocalPart)),
			fmt.Sprintf("We are going to recreate this resource if it is still part of your configuration, otherwise it will be removed from your state. Client error was: %v", err),
		)
		resp.State.RemoveResource(ctx)
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

	state.ID = types.StringValue(createMailboxID(state.DomainName, state.LocalPart))
	state.Address = types.StringValue(mailbox.Address)
	state.Name = types.StringValue(mailbox.Name)
	state.LocalPart = types.StringValue(mailbox.LocalPart)
	state.IsInternal = types.BoolValue(mailbox.IsInternal)
	state.MaySend = types.BoolValue(mailbox.MaySend)
	state.MayReceive = types.BoolValue(mailbox.MayReceive)
	state.MayAccessImap = types.BoolValue(mailbox.MayAccessImap)
	state.MayAccessPop3 = types.BoolValue(mailbox.MayAccessPop3)
	state.MayAccessManageSieve = types.BoolValue(mailbox.MayAccessManageSieve)
	if state.Password.IsUnknown() || state.Password.IsNull() {
		state.Password = types.StringValue("")
	}
	state.PasswordRecoveryEmail = types.StringValue(mailbox.PasswordRecoveryEmail)
	state.SpamAction = types.StringValue(mailbox.SpamAction)
	state.SpamAggressiveness = types.StringValue(mailbox.SpamAggressiveness)
	state.Expirable = types.BoolValue(mailbox.Expirable)
	state.ExpiresOn = types.StringValue(mailbox.ExpiresOn)
	state.RemoveUponExpiry = types.BoolValue(mailbox.RemoveUponExpiry)
	state.SenderDenyList = senderDenyList
	state.SenderDenyListPunycode = senderDenyListPunycode
	state.SenderAllowList = senderAllowList
	state.SenderAllowListPunycode = senderAllowListPunycode
	state.RecipientDenyList = recipientDenyList
	state.RecipientDenyListPunycode = recipientDenyListPunycode
	state.Delegations = delegations
	state.DelegationsPunycode = delegationsPunycode
	state.Identities = identities
	state.IdentitiesPunycode = identitiesPunycode
	state.AutoRespondActive = types.BoolValue(mailbox.AutoRespondActive)
	state.AutoRespondSubject = types.StringValue(mailbox.AutoRespondSubject)
	state.AutoRespondBody = types.StringValue(mailbox.AutoRespondBody)
	state.AutoRespondExpiresOn = types.StringValue(mailbox.AutoRespondExpiresOn)
	state.FooterActive = types.BoolValue(mailbox.FooterActive)
	state.FooterPlainBody = types.StringValue(mailbox.FooterPlainBody)
	state.FooterHtmlBody = types.StringValue(mailbox.FooterHtmlBody)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *mailboxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan mailboxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var wantedSenderDenyList []string
	if !plan.SenderDenyList.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderDenyList.ElementsAs(ctx, &wantedSenderDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.SenderDenyListPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderDenyListPunycode.ElementsAs(ctx, &wantedSenderDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var wantedSenderAllowList []string
	if !plan.SenderAllowList.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderAllowList.ElementsAs(ctx, &wantedSenderAllowList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.SenderAllowListPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderAllowListPunycode.ElementsAs(ctx, &wantedSenderAllowList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var wantedRecipientDenyList []string
	if !plan.RecipientDenyList.IsUnknown() {
		resp.Diagnostics.Append(plan.RecipientDenyList.ElementsAs(ctx, &wantedRecipientDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.RecipientDenyListPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.RecipientDenyListPunycode.ElementsAs(ctx, &wantedRecipientDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var wantedDelegations []string
	if !plan.Delegations.IsUnknown() {
		resp.Diagnostics.Append(plan.Delegations.ElementsAs(ctx, &wantedDelegations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.DelegationsPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.DelegationsPunycode.ElementsAs(ctx, &wantedDelegations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var wantedIdentities []string
	if !plan.Identities.IsUnknown() {
		resp.Diagnostics.Append(plan.Identities.ElementsAs(ctx, &wantedIdentities, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.IdentitiesPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.IdentitiesPunycode.ElementsAs(ctx, &wantedIdentities, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	mailbox := &model.Mailbox{
		Name:                  plan.Name.ValueString(),
		IsInternal:            plan.IsInternal.ValueBool(),
		MaySend:               plan.MaySend.ValueBool(),
		MayReceive:            plan.MayReceive.ValueBool(),
		MayAccessImap:         plan.MayAccessImap.ValueBool(),
		MayAccessPop3:         plan.MayAccessPop3.ValueBool(),
		MayAccessManageSieve:  plan.MayAccessManageSieve.ValueBool(),
		Password:              plan.Password.ValueString(),
		PasswordRecoveryEmail: plan.PasswordRecoveryEmail.ValueString(),
		SpamAction:            plan.SpamAction.ValueString(),
		SpamAggressiveness:    plan.SpamAggressiveness.ValueString(),
		Expirable:             plan.Expirable.ValueBool(),
		ExpiresOn:             plan.ExpiresOn.ValueString(),
		RemoveUponExpiry:      plan.RemoveUponExpiry.ValueBool(),
		SenderDenyList:        wantedSenderDenyList,
		SenderAllowList:       wantedSenderAllowList,
		RecipientDenyList:     wantedRecipientDenyList,
		Delegations:           wantedDelegations,
		Identities:            wantedIdentities,
		AutoRespondActive:     plan.AutoRespondActive.ValueBool(),
		AutoRespondSubject:    plan.AutoRespondSubject.ValueString(),
		AutoRespondBody:       plan.AutoRespondBody.ValueString(),
		AutoRespondExpiresOn:  plan.AutoRespondExpiresOn.ValueString(),
		FooterActive:          plan.FooterActive.ValueBool(),
		FooterPlainBody:       plan.FooterPlainBody.ValueString(),
		FooterHtmlBody:        plan.FooterHtmlBody.ValueString(),
	}

	if plan.Password.ValueString() != "" {
		mailbox.PasswordMethod = "password"
	} else if plan.PasswordRecoveryEmail.ValueString() != "" {
		mailbox.PasswordMethod = "invitation"
	}

	updatedMailbox, err := r.migaduClient.UpdateMailbox(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), mailbox)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating mailbox",
			fmt.Sprintf("Could not update mailbox %s: %v", createMailboxID(plan.DomainName, plan.LocalPart), err),
		)
		return
	}

	senderDenyList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(updatedMailbox.SenderDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderDenyListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(updatedMailbox.SenderDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderAllowList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(updatedMailbox.SenderAllowList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	senderAllowListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(updatedMailbox.SenderAllowList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	recipientDenyList, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(updatedMailbox.RecipientDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	recipientDenyListPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(updatedMailbox.RecipientDenyList, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	delegations, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(mailbox.Delegations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	delegationsPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(mailbox.Delegations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	identities, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(mailbox.Identities, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	identitiesPunycode, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToASCII(mailbox.Identities, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)

	plan.ID = types.StringValue(createMailboxID(plan.DomainName, plan.LocalPart))
	plan.Address = types.StringValue(updatedMailbox.Address)
	plan.Name = types.StringValue(updatedMailbox.Name)
	plan.LocalPart = types.StringValue(updatedMailbox.LocalPart)
	plan.IsInternal = types.BoolValue(updatedMailbox.IsInternal)
	plan.MaySend = types.BoolValue(updatedMailbox.MaySend)
	plan.MayReceive = types.BoolValue(updatedMailbox.MayReceive)
	plan.MayAccessImap = types.BoolValue(updatedMailbox.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(updatedMailbox.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(updatedMailbox.MayAccessManageSieve)
	plan.Password = types.StringValue(plan.Password.ValueString())
	plan.PasswordRecoveryEmail = types.StringValue(updatedMailbox.PasswordRecoveryEmail)
	plan.SpamAction = types.StringValue(updatedMailbox.SpamAction)
	plan.SpamAggressiveness = types.StringValue(updatedMailbox.SpamAggressiveness)
	plan.Expirable = types.BoolValue(updatedMailbox.Expirable)
	plan.ExpiresOn = types.StringValue(updatedMailbox.ExpiresOn)
	plan.RemoveUponExpiry = types.BoolValue(updatedMailbox.RemoveUponExpiry)
	plan.SenderDenyList = senderDenyList
	plan.SenderDenyListPunycode = senderDenyListPunycode
	plan.SenderAllowList = senderAllowList
	plan.SenderAllowListPunycode = senderAllowListPunycode
	plan.RecipientDenyList = recipientDenyList
	plan.RecipientDenyListPunycode = recipientDenyListPunycode
	plan.Delegations = delegations
	plan.DelegationsPunycode = delegationsPunycode
	plan.Identities = identities
	plan.IdentitiesPunycode = identitiesPunycode
	plan.AutoRespondActive = types.BoolValue(updatedMailbox.AutoRespondActive)
	plan.AutoRespondSubject = types.StringValue(updatedMailbox.AutoRespondSubject)
	plan.AutoRespondBody = types.StringValue(updatedMailbox.AutoRespondBody)
	plan.AutoRespondExpiresOn = types.StringValue(updatedMailbox.AutoRespondExpiresOn)
	plan.FooterActive = types.BoolValue(updatedMailbox.FooterActive)
	plan.FooterPlainBody = types.StringValue(updatedMailbox.FooterPlainBody)
	plan.FooterHtmlBody = types.StringValue(updatedMailbox.FooterHtmlBody)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *mailboxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state mailboxResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.migaduClient.DeleteMailbox(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting mailbox",
			fmt.Sprintf("Could not delete mailbox %s: %v", createMailboxID(state.DomainName, state.LocalPart), err),
		)
		return
	}
}

func (r *mailboxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "@")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Error importing mailbox",
			fmt.Sprintf("Expected import identifier with format: 'local_part@domain_name' Got: '%q'", req.ID),
		)
		return
	}

	localPart := idParts[0]
	domainName := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"local_part":  localPart,
		"domain_name": domainName,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("local_part"), localPart)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_name"), domainName)...)
}

func createMailboxID(domainName, localPart types.String) string {
	return fmt.Sprintf("%s@%s", localPart.ValueString(), domainName.ValueString())
}
