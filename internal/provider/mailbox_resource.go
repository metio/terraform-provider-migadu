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
	PasswordMethod            types.String  `tfsdk:"password_method"`
	Password                  types.String  `tfsdk:"password"`
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
	Identities                types.List    `tfsdk:"identities"`
}

func (r *mailboxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mailbox"
}

func (r *mailboxResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage a single mailbox.",
		MarkdownDescription: "Manage a single mailbox.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the mailbox to manage.",
				MarkdownDescription: "The domain name of the mailbox to manage.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the mailbox to manage.",
				MarkdownDescription: "The local part of the mailbox to manage.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Contains the full email address 'local_part@domain_name'.",
				MarkdownDescription: "Contains the full email address 'local_part@domain_name'.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				Description:         "Contains the full email address 'local_part@domain_name' as returned by the Migadu API. This might be different from the 'id' attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "Contains the full email address `local_part@domain_name` as returned by the Migadu API. This might be different from the `id` attribute in case you are using international domain names. The Migadu API always returns the punycode version of a domain.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"is_internal": schema.BoolAttribute{
				Description:         "Internal mailboxes can only receive emails from Migadu email servers.",
				MarkdownDescription: "Internal mailboxes can only receive emails from Migadu email servers.",
				Optional:            true,
				Computed:            true,
			},
			"may_send": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"may_receive": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"may_access_imap": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"may_access_pop3": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"may_access_manage_sieve": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"password_method": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"password_recovery_email": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"spam_action": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"spam_aggressiveness": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"expirable": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"expires_on": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"remove_upon_expiry": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"sender_denylist": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"sender_denylist_punycode": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"sender_allowlist": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"sender_allowlist_punycode": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"recipient_denylist": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"recipient_denylist_punycode": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"auto_respond_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_respond_subject": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_respond_body": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"auto_respond_expires_on": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"footer_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"footer_plain_body": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"footer_html_body": schema.StringAttribute{
				Optional: true,
				Computed: true,
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

	var senderDenyList []string
	if !plan.SenderDenyList.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderDenyList.ElementsAs(ctx, &senderDenyList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.SenderDenyListPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.SenderDenyListPunycode.ElementsAs(ctx, &senderDenyList, false)...)
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
		PasswordMethod:        plan.PasswordMethod.ValueString(),
		Password:              plan.Password.ValueString(),
		PasswordRecoveryEmail: plan.PasswordRecoveryEmail.ValueString(),
		SpamAction:            plan.SpamAction.ValueString(),
		SpamAggressiveness:    plan.SpamAggressiveness.ValueString(),
		Expirable:             plan.Expirable.ValueBool(),
		ExpiresOn:             plan.ExpiresOn.ValueString(),
		RemoveUponExpiry:      plan.RemoveUponExpiry.ValueBool(),
		SenderDenyList:        senderDenyList,
		SenderAllowList:       nil,
		RecipientDenyList:     nil,
		AutoRespondActive:     plan.AutoRespondActive.ValueBool(),
		AutoRespondSubject:    plan.AutoRespondSubject.ValueString(),
		AutoRespondBody:       plan.AutoRespondBody.ValueString(),
		AutoRespondExpiresOn:  plan.AutoRespondExpiresOn.ValueString(),
		FooterActive:          plan.FooterActive.ValueBool(),
		FooterPlainBody:       plan.FooterPlainBody.ValueString(),
		FooterHtmlBody:        plan.FooterHtmlBody.ValueString(),
	}

	createdMailbox, err := r.migaduClient.CreateMailbox(ctx, plan.DomainName.ValueString(), mailbox)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating mailbox",
			fmt.Sprintf("Could not create mailbox %s: %v", createMailboxID(plan.DomainName, plan.LocalPart), err),
		)
		return
	}

	plan.ID = types.StringValue(createMailboxID(plan.DomainName, plan.LocalPart))
	plan.Name = types.StringValue(createdMailbox.Name)
	plan.LocalPart = types.StringValue(createdMailbox.LocalPart)
	plan.IsInternal = types.BoolValue(createdMailbox.IsInternal)
	plan.MaySend = types.BoolValue(createdMailbox.MaySend)
	plan.MayReceive = types.BoolValue(createdMailbox.MayReceive)
	plan.MayAccessImap = types.BoolValue(createdMailbox.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(createdMailbox.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(createdMailbox.MayAccessManageSieve)
	plan.PasswordMethod = types.StringValue(createdMailbox.PasswordMethod)
	plan.Password = types.StringValue(createdMailbox.Password)
	plan.PasswordRecoveryEmail = types.StringValue(createdMailbox.PasswordRecoveryEmail)
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
		resp.Diagnostics.AddError(
			"Error reading mailbox",
			fmt.Sprintf("Could not read mailbox %s: %v", createMailboxID(state.DomainName, state.LocalPart), err),
		)
		return
	}

	state.ID = types.StringValue(createMailboxID(state.DomainName, state.LocalPart))
	state.Name = types.StringValue(mailbox.Name)
	state.LocalPart = types.StringValue(mailbox.LocalPart)
	state.IsInternal = types.BoolValue(mailbox.IsInternal)
	state.MaySend = types.BoolValue(mailbox.MaySend)
	state.MayReceive = types.BoolValue(mailbox.MayReceive)
	state.MayAccessImap = types.BoolValue(mailbox.MayAccessImap)
	state.MayAccessPop3 = types.BoolValue(mailbox.MayAccessPop3)
	state.MayAccessManageSieve = types.BoolValue(mailbox.MayAccessManageSieve)
	state.PasswordMethod = types.StringValue(mailbox.PasswordMethod)
	state.Password = types.StringValue(mailbox.Password)
	state.PasswordRecoveryEmail = types.StringValue(mailbox.PasswordRecoveryEmail)
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

	mailbox := &model.Mailbox{
		Name:                  plan.Name.ValueString(),
		IsInternal:            plan.IsInternal.ValueBool(),
		MaySend:               plan.MaySend.ValueBool(),
		MayReceive:            plan.MayReceive.ValueBool(),
		MayAccessImap:         plan.MayAccessImap.ValueBool(),
		MayAccessPop3:         plan.MayAccessPop3.ValueBool(),
		MayAccessManageSieve:  plan.MayAccessManageSieve.ValueBool(),
		PasswordMethod:        plan.PasswordMethod.ValueString(),
		Password:              plan.Password.ValueString(),
		PasswordRecoveryEmail: plan.PasswordRecoveryEmail.ValueString(),
		SpamAction:            plan.SpamAction.ValueString(),
		SpamAggressiveness:    plan.SpamAggressiveness.ValueString(),
		Expirable:             plan.Expirable.ValueBool(),
		ExpiresOn:             plan.ExpiresOn.ValueString(),
		RemoveUponExpiry:      plan.RemoveUponExpiry.ValueBool(),
		SenderDenyList:        nil,
		SenderAllowList:       nil,
		RecipientDenyList:     nil,
		AutoRespondActive:     plan.AutoRespondActive.ValueBool(),
		AutoRespondSubject:    plan.AutoRespondSubject.ValueString(),
		AutoRespondBody:       plan.AutoRespondBody.ValueString(),
		AutoRespondExpiresOn:  plan.AutoRespondExpiresOn.ValueString(),
		FooterActive:          plan.FooterActive.ValueBool(),
		FooterPlainBody:       plan.FooterPlainBody.ValueString(),
		FooterHtmlBody:        plan.FooterHtmlBody.ValueString(),
	}

	updatedMailbox, err := r.migaduClient.UpdateMailbox(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), mailbox)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating mailbox",
			fmt.Sprintf("Could not update mailbox %s: %v", createMailboxID(plan.DomainName, plan.LocalPart), err),
		)
		return
	}

	plan.ID = types.StringValue(createMailboxID(plan.DomainName, plan.LocalPart))
	plan.Name = types.StringValue(updatedMailbox.Name)
	plan.LocalPart = types.StringValue(updatedMailbox.LocalPart)
	plan.IsInternal = types.BoolValue(updatedMailbox.IsInternal)
	plan.MaySend = types.BoolValue(updatedMailbox.MaySend)
	plan.MayReceive = types.BoolValue(updatedMailbox.MayReceive)
	plan.MayAccessImap = types.BoolValue(updatedMailbox.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(updatedMailbox.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(updatedMailbox.MayAccessManageSieve)
	plan.PasswordMethod = types.StringValue(updatedMailbox.PasswordMethod)
	plan.Password = types.StringValue(updatedMailbox.Password)
	plan.PasswordRecoveryEmail = types.StringValue(updatedMailbox.PasswordRecoveryEmail)
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
