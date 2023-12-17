/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/migadu-client.go/client"
	"github.com/metio/migadu-client.go/model"
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
	"net/http"
	"strings"
)

var (
	_ resource.Resource                = (*MailboxResource)(nil)
	_ resource.ResourceWithConfigure   = (*MailboxResource)(nil)
	_ resource.ResourceWithImportState = (*MailboxResource)(nil)
)

func NewMailboxResource() resource.Resource {
	return &MailboxResource{}
}

type MailboxResource struct {
	MigaduClient *client.MigaduClient
}

type MailboxResourceModel struct {
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
	Password              types.String                      `tfsdk:"password"`
	PasswordRecoveryEmail custom_types.EmailAddressValue    `tfsdk:"password_recovery_email"`
	PasswordMethod        types.String                      `tfsdk:"password_method"`
	SpamAction            types.String                      `tfsdk:"spam_action"`
	SpamAggressiveness    types.String                      `tfsdk:"spam_aggressiveness"`
	Expirable             types.Bool                        `tfsdk:"expirable"`
	ExpiresOn             types.String                      `tfsdk:"expires_on"`
	RemoveUponExpiry      types.Bool                        `tfsdk:"remove_upon_expiry"`
	SenderDenyList        custom_types.EmailAddressSetValue `tfsdk:"sender_denylist"`
	SenderAllowList       custom_types.EmailAddressSetValue `tfsdk:"sender_allowlist"`
	RecipientDenyList     custom_types.EmailAddressSetValue `tfsdk:"recipient_denylist"`
	Delegations           custom_types.EmailAddressSetValue `tfsdk:"delegations"`
	Identities            custom_types.EmailAddressSetValue `tfsdk:"identities"`
	AutoRespondActive     types.Bool                        `tfsdk:"auto_respond_active"`
	AutoRespondSubject    types.String                      `tfsdk:"auto_respond_subject"`
	AutoRespondBody       types.String                      `tfsdk:"auto_respond_body"`
	AutoRespondExpiresOn  types.String                      `tfsdk:"auto_respond_expires_on"`
	FooterActive          types.Bool                        `tfsdk:"footer_active"`
	FooterPlainBody       types.String                      `tfsdk:"footer_plain_body"`
	FooterHtmlBody        types.String                      `tfsdk:"footer_html_body"`
}

func (r *MailboxResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mailbox"
}

func (r *MailboxResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Provides a mailbox.",
		MarkdownDescription: "Provides a mailbox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Contains the value 'local_part@domain_name'.",
				MarkdownDescription: "Contains the value `local_part@domain_name`.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Description:         "The email address of the mailbox 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "The email address of the mailbox `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the mailbox.",
				MarkdownDescription: "The name of the mailbox.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"is_internal": schema.BoolAttribute{
				Description:         "Whether this mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
				MarkdownDescription: "Whether this mailbox is internal only. An internal mailbox can only receive emails from Migadu servers.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_send": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to send emails.",
				MarkdownDescription: "Whether this mailbox is allowed to send emails.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_receive": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to receive emails.",
				MarkdownDescription: "Whether this mailbox is allowed to receive emails.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_access_imap": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to use IMAP.",
				MarkdownDescription: "Whether this mailbox is allowed to use IMAP.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_access_pop3": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to use POP3.",
				MarkdownDescription: "Whether this mailbox is allowed to use POP3.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_access_manage_sieve": schema.BoolAttribute{
				Description:         "Whether this mailbox is allowed to manage the mail sieve.",
				MarkdownDescription: "Whether this mailbox is allowed to manage the mail sieve.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				Description:         "The password of this mailbox.",
				MarkdownDescription: "The password of this mailbox.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("password_recovery_email")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password_recovery_email": schema.StringAttribute{
				Description:         "The recovery email address of this mailbox.",
				MarkdownDescription: "The recovery email address of this mailbox.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("password")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password_method": schema.StringAttribute{
				Description:         "The password method of this mailbox. If this is set to 'invitation' an email will be send to the 'password_recovery_email' and users can set their own password.",
				MarkdownDescription: "The password method of this mailbox. If this is set to 'invitation' an email will be send to the 'password_recovery_email' and users can set their own password.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("password", "invitation"),
				},
				Default: stringdefault.StaticString("password"),
			},
			"spam_action": schema.StringAttribute{
				Description:         "The action to take once spam arrives in this mailbox.",
				MarkdownDescription: "The action to take once spam arrives in this mailbox.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"spam_aggressiveness": schema.StringAttribute{
				Description:         "How aggressive will spam be detected in this mailbox.",
				MarkdownDescription: "How aggressive will spam be detected in this mailbox.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"expirable": schema.BoolAttribute{
				Description:         "Whether this mailbox expires in the future.",
				MarkdownDescription: "Whether this mailbox expires in the future.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"expires_on": schema.StringAttribute{
				Description:         "The expiration date of this mailbox.",
				MarkdownDescription: "The expiration date of this mailbox.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Description:         "Whether this mailbox will be removed upon expiry.",
				MarkdownDescription: "Whether this mailbox will be removed upon expiry.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"sender_denylist": schema.SetAttribute{
				Description:         "The email addresses of senders that will always be denied delivery.",
				MarkdownDescription: "The email addresses of senders that will always be denied delivery.",
				Required:            false,
				Optional:            true,
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
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
				CustomType: custom_types.EmailAddressSetType{
					SetType: types.SetType{
						ElemType: custom_types.EmailAddressType{},
					},
				},
			},
			"delegations": schema.SetAttribute{
				Description:         "The delegations of the mailbox.",
				MarkdownDescription: "The delegations of the mailbox.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				CustomType: custom_types.EmailAddressSetType{
					SetType: types.SetType{
						ElemType: custom_types.EmailAddressType{},
					},
				},
			},
			"identities": schema.SetAttribute{
				Description:         "The identities of the mailbox.",
				MarkdownDescription: "The identities of the mailbox.",
				Required:            false,
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
			"auto_respond_subject": schema.StringAttribute{
				Description:         "The subject of the automatic response.",
				MarkdownDescription: "The subject of the automatic response.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"auto_respond_body": schema.StringAttribute{
				Description:         "The body of the automatic response.",
				MarkdownDescription: "The body of the automatic response.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"auto_respond_expires_on": schema.StringAttribute{
				Description:         "The expiration date of the automatic response.",
				MarkdownDescription: "The expiration date of the automatic response.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"footer_active": schema.BoolAttribute{
				Description:         "Whether the footer of this mailbox is active.",
				MarkdownDescription: "Whether the footer of this mailbox is active.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"footer_plain_body": schema.StringAttribute{
				Description:         "The footer of this mailbox in text/plain format.",
				MarkdownDescription: "The footer of this mailbox in text/plain format.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"footer_html_body": schema.StringAttribute{
				Description:         "The footer of this mailbox in text/html format.",
				MarkdownDescription: "The footer of this mailbox in text/html format.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}
func (r *MailboxResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	if migaduClient, ok := request.ProviderData.(*client.MigaduClient); ok {
		r.MigaduClient = migaduClient
	} else {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.MigaduClient, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
}

func (r *MailboxResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan MailboxResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	if plan.PasswordMethod.ValueString() == "password" && plan.Password.ValueString() == "" {
		response.Diagnostics.AddError(
			"Error creating mailbox",
			"Cannot use 'password_method = password' without a 'password'",
		)
		return
	}
	if plan.PasswordMethod.ValueString() == "invitation" && plan.PasswordRecoveryEmail.ValueString() == "" {
		response.Diagnostics.AddError(
			"Error creating mailbox",
			"Cannot use 'password_method = invitation' without a 'password_recovery_email'",
		)
		return
	}

	var senderDenyList []string
	if plan.SenderDenyList.IsUnknown() {
		plan.SenderDenyList = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.SenderDenyList.ElementsAs(ctx, &senderDenyList, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	var senderAllowList []string
	if plan.SenderAllowList.IsUnknown() {
		plan.SenderAllowList = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.SenderAllowList.ElementsAs(ctx, &senderAllowList, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	var recipientDenyList []string
	if plan.RecipientDenyList.IsUnknown() {
		plan.RecipientDenyList = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.RecipientDenyList.ElementsAs(ctx, &recipientDenyList, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	var delegations []string
	if plan.Delegations.IsUnknown() {
		plan.Delegations = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.Delegations.ElementsAs(ctx, &delegations, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	var identities []string
	if plan.Identities.IsUnknown() {
		plan.Identities = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.Identities.ElementsAs(ctx, &identities, false)...)
		if response.Diagnostics.HasError() {
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
		PasswordMethod:        plan.PasswordMethod.ValueString(),
		SpamAction:            plan.SpamAction.ValueString(),
		SpamAggressiveness:    plan.SpamAggressiveness.ValueString(),
		Expirable:             plan.Expirable.ValueBool(),
		ExpiresOn:             plan.ExpiresOn.ValueString(),
		RemoveUponExpiry:      plan.RemoveUponExpiry.ValueBool(),
		SenderDenyList:        senderDenyList,
		SenderAllowList:       senderAllowList,
		RecipientDenyList:     recipientDenyList,
		Delegations:           delegations,
		Identities:            identities,
		AutoRespondActive:     plan.AutoRespondActive.ValueBool(),
		AutoRespondSubject:    plan.AutoRespondSubject.ValueString(),
		AutoRespondBody:       plan.AutoRespondBody.ValueString(),
		AutoRespondExpiresOn:  plan.AutoRespondExpiresOn.ValueString(),
		FooterActive:          plan.FooterActive.ValueBool(),
		FooterPlainBody:       plan.FooterPlainBody.ValueString(),
		FooterHtmlBody:        plan.FooterHtmlBody.ValueString(),
	}

	createdMailbox, err := r.MigaduClient.CreateMailbox(ctx, plan.DomainName.ValueString(), mailbox)
	if err != nil {
		response.Diagnostics.Append(MailboxCreateError(err))
		return
	}

	plan.ID = custom_types.NewEmailAddressValue(CreateMailboxID(plan.LocalPart, plan.DomainName))
	plan.Address = custom_types.NewEmailAddressValue(createdMailbox.Address)
	plan.Name = types.StringValue(createdMailbox.Name)
	plan.IsInternal = types.BoolValue(createdMailbox.IsInternal)
	plan.MaySend = types.BoolValue(createdMailbox.MaySend)
	plan.MayReceive = types.BoolValue(createdMailbox.MayReceive)
	plan.MayAccessImap = types.BoolValue(createdMailbox.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(createdMailbox.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(createdMailbox.MayAccessManageSieve)
	plan.Password = types.StringValue(plan.Password.ValueString())
	plan.PasswordRecoveryEmail = custom_types.NewEmailAddressValue(plan.PasswordRecoveryEmail.ValueString())
	plan.SpamAction = types.StringValue(createdMailbox.SpamAction)
	plan.SpamAggressiveness = types.StringValue(createdMailbox.SpamAggressiveness)
	plan.Expirable = types.BoolValue(createdMailbox.Expirable)
	plan.ExpiresOn = types.StringValue(createdMailbox.ExpiresOn)
	plan.RemoveUponExpiry = types.BoolValue(createdMailbox.RemoveUponExpiry)
	plan.AutoRespondActive = types.BoolValue(createdMailbox.AutoRespondActive)
	plan.AutoRespondSubject = types.StringValue(createdMailbox.AutoRespondSubject)
	plan.AutoRespondBody = types.StringValue(createdMailbox.AutoRespondBody)
	plan.AutoRespondExpiresOn = types.StringValue(createdMailbox.AutoRespondExpiresOn)
	plan.FooterActive = types.BoolValue(createdMailbox.FooterActive)
	plan.FooterPlainBody = types.StringValue(createdMailbox.FooterPlainBody)
	plan.FooterHtmlBody = types.StringValue(createdMailbox.FooterHtmlBody)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (r *MailboxResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state MailboxResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	mailbox, err := r.MigaduClient.GetMailbox(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString())
	if err != nil {
		var requestError *client.RequestError
		if errors.As(err, &requestError) {
			if requestError.StatusCode == http.StatusNotFound {
				response.State.RemoveResource(ctx)
				return
			}
		}
		response.Diagnostics.Append(MailboxReadError(err))
		return
	}

	senderDenyList, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.SenderDenyList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	senderDenyListEqual, diags := state.SenderDenyList.SetSemanticEquals(ctx, senderDenyList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if !senderDenyListEqual {
		state.SenderDenyList = senderDenyList
	}

	senderAllowList, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.SenderAllowList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	senderAllowListEqual, diags := state.SenderAllowList.SetSemanticEquals(ctx, senderAllowList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if !senderAllowListEqual {
		state.SenderAllowList = senderAllowList
	}

	recipientDenyList, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.RecipientDenyList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	recipientDenyListEqual, diags := state.RecipientDenyList.SetSemanticEquals(ctx, recipientDenyList)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if !recipientDenyListEqual {
		state.RecipientDenyList = recipientDenyList
	}

	delegations, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.Delegations)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	delegationsEqual, diags := state.Delegations.SetSemanticEquals(ctx, delegations)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if !delegationsEqual {
		state.Delegations = delegations
	}

	identities, diags := custom_types.NewEmailAddressSetValueFrom(ctx, mailbox.Identities)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	identitiesEqual, diags := state.Identities.SetSemanticEquals(ctx, identities)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if !identitiesEqual {
		state.Identities = identities
	}

	state.ID = custom_types.NewEmailAddressValue(CreateMailboxID(state.LocalPart, state.DomainName))
	state.Address = custom_types.NewEmailAddressValue(mailbox.Address)
	state.Name = types.StringValue(mailbox.Name)
	state.LocalPart = types.StringValue(mailbox.LocalPart)
	state.IsInternal = types.BoolValue(mailbox.IsInternal)
	state.MaySend = types.BoolValue(mailbox.MaySend)
	state.MayReceive = types.BoolValue(mailbox.MayReceive)
	state.MayAccessImap = types.BoolValue(mailbox.MayAccessImap)
	state.MayAccessPop3 = types.BoolValue(mailbox.MayAccessPop3)
	state.MayAccessManageSieve = types.BoolValue(mailbox.MayAccessManageSieve)
	if state.Password.IsUnknown() {
		state.Password = types.StringNull()
	}
	state.PasswordRecoveryEmail = custom_types.NewEmailAddressValue(mailbox.PasswordRecoveryEmail)
	if state.PasswordMethod.IsUnknown() || state.PasswordMethod.IsNull() {
		state.PasswordMethod = types.StringValue("password")
	}
	state.SpamAction = types.StringValue(mailbox.SpamAction)
	state.SpamAggressiveness = types.StringValue(mailbox.SpamAggressiveness)
	state.Expirable = types.BoolValue(mailbox.Expirable)
	state.ExpiresOn = types.StringValue(mailbox.ExpiresOn)
	state.RemoveUponExpiry = types.BoolValue(mailbox.RemoveUponExpiry)
	state.AutoRespondActive = types.BoolValue(mailbox.AutoRespondActive)
	state.AutoRespondSubject = types.StringValue(mailbox.AutoRespondSubject)
	state.AutoRespondBody = types.StringValue(mailbox.AutoRespondBody)
	state.AutoRespondExpiresOn = types.StringValue(mailbox.AutoRespondExpiresOn)
	state.FooterActive = types.BoolValue(mailbox.FooterActive)
	state.FooterPlainBody = types.StringValue(mailbox.FooterPlainBody)
	state.FooterHtmlBody = types.StringValue(mailbox.FooterHtmlBody)

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *MailboxResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan MailboxResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var senderDenyList []string
	if plan.SenderDenyList.IsUnknown() {
		plan.SenderDenyList = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.SenderDenyList.ElementsAs(ctx, &senderDenyList, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	var senderAllowList []string
	if plan.SenderAllowList.IsUnknown() {
		plan.SenderAllowList = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.SenderAllowList.ElementsAs(ctx, &senderAllowList, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	var recipientDenyList []string
	if plan.RecipientDenyList.IsUnknown() {
		plan.RecipientDenyList = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.RecipientDenyList.ElementsAs(ctx, &recipientDenyList, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	var delegations []string
	if plan.Delegations.IsUnknown() {
		plan.Delegations = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.Delegations.ElementsAs(ctx, &delegations, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	var identities []string
	if plan.Identities.IsUnknown() {
		plan.Identities = custom_types.NewEmailAddressSetNull()
	} else {
		response.Diagnostics.Append(plan.Identities.ElementsAs(ctx, &identities, false)...)
		if response.Diagnostics.HasError() {
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
		SenderDenyList:        senderDenyList,
		SenderAllowList:       senderAllowList,
		RecipientDenyList:     recipientDenyList,
		Delegations:           delegations,
		Identities:            identities,
		AutoRespondActive:     plan.AutoRespondActive.ValueBool(),
		AutoRespondSubject:    plan.AutoRespondSubject.ValueString(),
		AutoRespondBody:       plan.AutoRespondBody.ValueString(),
		AutoRespondExpiresOn:  plan.AutoRespondExpiresOn.ValueString(),
		FooterActive:          plan.FooterActive.ValueBool(),
		FooterPlainBody:       plan.FooterPlainBody.ValueString(),
		FooterHtmlBody:        plan.FooterHtmlBody.ValueString(),
	}

	updatedMailbox, err := r.MigaduClient.UpdateMailbox(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), mailbox)
	if err != nil {
		response.Diagnostics.Append(MailboxUpdateError(err))
		return
	}

	plan.ID = custom_types.NewEmailAddressValue(CreateMailboxID(plan.LocalPart, plan.DomainName))
	plan.Address = custom_types.NewEmailAddressValue(updatedMailbox.Address)
	plan.Name = types.StringValue(updatedMailbox.Name)
	plan.IsInternal = types.BoolValue(updatedMailbox.IsInternal)
	plan.MaySend = types.BoolValue(updatedMailbox.MaySend)
	plan.MayReceive = types.BoolValue(updatedMailbox.MayReceive)
	plan.MayAccessImap = types.BoolValue(updatedMailbox.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(updatedMailbox.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(updatedMailbox.MayAccessManageSieve)
	plan.Password = types.StringValue(plan.Password.ValueString())
	plan.PasswordRecoveryEmail = custom_types.NewEmailAddressValue(plan.PasswordRecoveryEmail.ValueString())
	plan.SpamAction = types.StringValue(updatedMailbox.SpamAction)
	plan.SpamAggressiveness = types.StringValue(updatedMailbox.SpamAggressiveness)
	plan.Expirable = types.BoolValue(updatedMailbox.Expirable)
	plan.ExpiresOn = types.StringValue(updatedMailbox.ExpiresOn)
	plan.RemoveUponExpiry = types.BoolValue(updatedMailbox.RemoveUponExpiry)
	plan.AutoRespondActive = types.BoolValue(updatedMailbox.AutoRespondActive)
	plan.AutoRespondSubject = types.StringValue(updatedMailbox.AutoRespondSubject)
	plan.AutoRespondBody = types.StringValue(updatedMailbox.AutoRespondBody)
	plan.AutoRespondExpiresOn = types.StringValue(updatedMailbox.AutoRespondExpiresOn)
	plan.FooterActive = types.BoolValue(updatedMailbox.FooterActive)
	plan.FooterPlainBody = types.StringValue(updatedMailbox.FooterPlainBody)
	plan.FooterHtmlBody = types.StringValue(updatedMailbox.FooterHtmlBody)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (r *MailboxResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state MailboxResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := r.MigaduClient.DeleteMailbox(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString())
	if err != nil {
		response.Diagnostics.Append(MailboxDeleteError(err))
		return
	}
}

func (r *MailboxResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	idParts := strings.Split(request.ID, "@")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.Append(MailboxImportError(request.ID))
		return
	}

	localPart := idParts[0]
	domainName := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"local_part":  localPart,
		"domain_name": domainName,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("local_part"), localPart)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("domain_name"), domainName)...)
}
