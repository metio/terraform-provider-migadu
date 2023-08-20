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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-migadu/internal/provider/custom_types"
	"github.com/metio/terraform-provider-migadu/migadu/client"
	"github.com/metio/terraform-provider-migadu/migadu/model"
	"net/http"
	"strings"
)

var (
	_ resource.Resource                = (*IdentityResource)(nil)
	_ resource.ResourceWithConfigure   = (*IdentityResource)(nil)
	_ resource.ResourceWithImportState = (*IdentityResource)(nil)
)

func NewIdentityResource() resource.Resource {
	return &IdentityResource{}
}

type IdentityResource struct {
	MigaduClient *client.MigaduClient
}

type IdentityResourceModel struct {
	ID                   types.String                   `tfsdk:"id"`
	LocalPart            types.String                   `tfsdk:"local_part"`
	DomainName           custom_types.DomainNameValue   `tfsdk:"domain_name"`
	Identity             types.String                   `tfsdk:"identity"`
	Address              custom_types.EmailAddressValue `tfsdk:"address"`
	Name                 types.String                   `tfsdk:"name"`
	MaySend              types.Bool                     `tfsdk:"may_send"`
	MayReceive           types.Bool                     `tfsdk:"may_receive"`
	MayAccessImap        types.Bool                     `tfsdk:"may_access_imap"`
	MayAccessPop3        types.Bool                     `tfsdk:"may_access_pop3"`
	MayAccessManageSieve types.Bool                     `tfsdk:"may_access_manage_sieve"`
	Password             types.String                   `tfsdk:"password"`
	FooterActive         types.Bool                     `tfsdk:"footer_active"`
	FooterPlainBody      types.String                   `tfsdk:"footer_plain_body"`
	FooterHtmlBody       types.String                   `tfsdk:"footer_html_body"`
}

func (r *IdentityResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_identity"
}

func (r *IdentityResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Provides an identity to an existing mailbox.",
		MarkdownDescription: "Provides an identity to an existing mailbox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Contains the value 'local_part@domain_name/identity'.",
				MarkdownDescription: "Contains the value `local_part@domain_name/identity`.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the mailbox that owns the identity.",
				MarkdownDescription: "The local part of the mailbox that owns the identity.",
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
				Description:         "The domain name of the mailbox/identity.",
				MarkdownDescription: "The domain name of the mailbox/identity.",
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
			"identity": schema.StringAttribute{
				Description:         "The local part of the identity.",
				MarkdownDescription: "The local part of the identity.",
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
			"address": schema.StringAttribute{
				Description:         "Contains the email address of the identity 'identity@domain_name' as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "Contains the email address of the identity `identity@domain_name` as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				CustomType:          custom_types.EmailAddressType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the identity.",
				MarkdownDescription: "The name of the identity.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_send": schema.BoolAttribute{
				Description:         "Whether the identity is allowed to send emails.",
				MarkdownDescription: "Whether the identity is allowed to send emails.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_receive": schema.BoolAttribute{
				Description:         "Whether the identity is allowed to receive emails.",
				MarkdownDescription: "Whether the identity is allowed to receive emails.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_access_imap": schema.BoolAttribute{
				Description:         "Whether the identity is allowed to use IMAP.",
				MarkdownDescription: "Whether the identity is allowed to use IMAP.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_access_pop3": schema.BoolAttribute{
				Description:         "Whether the identity is allowed to use POP3.",
				MarkdownDescription: "Whether the identity is allowed to use POP3.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"may_access_manage_sieve": schema.BoolAttribute{
				Description:         "Whether the identity is allowed to manage the mail sieve.",
				MarkdownDescription: "Whether the identity is allowed to manage the mail sieve.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				Description:         "The password of the identity.",
				MarkdownDescription: "The password of the identity.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"footer_active": schema.BoolAttribute{
				Description:         "Whether the footer of the identity is active.",
				MarkdownDescription: "Whether the footer of the identity is active.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"footer_plain_body": schema.StringAttribute{
				Description:         "The footer of the identity in 'text/plain' format.",
				MarkdownDescription: "The footer of the identity in `text/plain` format.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
			"footer_html_body": schema.StringAttribute{
				Description:         "The footer of the identity in 'text/html' format.",
				MarkdownDescription: "The footer of the identity in `text/html` format.",
				Required:            false,
				Optional:            true,
				Computed:            true,
			},
		},
	}
}
func (r *IdentityResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *IdentityResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan IdentityResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	identity := &model.Identity{
		LocalPart:            plan.Identity.ValueString(),
		Name:                 plan.Name.ValueString(),
		MaySend:              plan.MaySend.ValueBool(),
		MayReceive:           plan.MayReceive.ValueBool(),
		MayAccessImap:        plan.MayAccessImap.ValueBool(),
		MayAccessPop3:        plan.MayAccessPop3.ValueBool(),
		MayAccessManageSieve: plan.MayAccessManageSieve.ValueBool(),
		Password:             plan.Password.ValueString(),
		FooterActive:         plan.FooterActive.ValueBool(),
		FooterPlainBody:      plan.FooterPlainBody.ValueString(),
		FooterHtmlBody:       plan.FooterHtmlBody.ValueString(),
	}

	createdIdentity, err := r.MigaduClient.CreateIdentity(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), identity)
	if err != nil {
		response.Diagnostics.Append(IdentityCreateError(err))
		return
	}

	plan.ID = types.StringValue(CreateIdentityID(plan.LocalPart, plan.DomainName, plan.Identity))
	plan.Address = custom_types.NewEmailAddressValue(createdIdentity.Address)
	plan.Name = types.StringValue(createdIdentity.Name)
	plan.MaySend = types.BoolValue(createdIdentity.MaySend)
	plan.MayReceive = types.BoolValue(createdIdentity.MayReceive)
	plan.MayAccessImap = types.BoolValue(createdIdentity.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(createdIdentity.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(createdIdentity.MayAccessManageSieve)
	plan.FooterActive = types.BoolValue(createdIdentity.FooterActive)
	plan.FooterPlainBody = types.StringValue(createdIdentity.FooterPlainBody)
	plan.FooterHtmlBody = types.StringValue(createdIdentity.FooterHtmlBody)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (r *IdentityResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state IdentityResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	identity, err := r.MigaduClient.GetIdentity(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString(), state.Identity.ValueString())
	if err != nil {
		var requestError *client.RequestError
		if errors.As(err, &requestError) {
			if requestError.StatusCode == http.StatusNotFound {
				response.State.RemoveResource(ctx)
				return
			}
		}
		response.Diagnostics.Append(IdentityReadError(err))
		return
	}

	state.ID = types.StringValue(CreateIdentityID(state.LocalPart, state.DomainName, state.Identity))
	state.Address = custom_types.NewEmailAddressValue(identity.Address)
	state.Name = types.StringValue(identity.Name)
	state.MaySend = types.BoolValue(identity.MaySend)
	state.MayReceive = types.BoolValue(identity.MayReceive)
	state.MayAccessImap = types.BoolValue(identity.MayAccessImap)
	state.MayAccessPop3 = types.BoolValue(identity.MayAccessPop3)
	state.MayAccessManageSieve = types.BoolValue(identity.MayAccessManageSieve)
	state.FooterActive = types.BoolValue(identity.FooterActive)
	state.FooterPlainBody = types.StringValue(identity.FooterPlainBody)
	state.FooterHtmlBody = types.StringValue(identity.FooterHtmlBody)

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *IdentityResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan IdentityResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	identity := &model.Identity{
		Name:                 plan.Name.ValueString(),
		MaySend:              plan.MaySend.ValueBool(),
		MayReceive:           plan.MayReceive.ValueBool(),
		MayAccessImap:        plan.MayAccessImap.ValueBool(),
		MayAccessPop3:        plan.MayAccessPop3.ValueBool(),
		MayAccessManageSieve: plan.MayAccessManageSieve.ValueBool(),
		Password:             plan.Password.ValueString(),
		FooterActive:         plan.FooterActive.ValueBool(),
		FooterPlainBody:      plan.FooterPlainBody.ValueString(),
		FooterHtmlBody:       plan.FooterHtmlBody.ValueString(),
	}

	updatedIdentity, err := r.MigaduClient.UpdateIdentity(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), plan.Identity.ValueString(), identity)
	if err != nil {
		response.Diagnostics.Append(IdentityUpdateError(err))
		return
	}

	plan.ID = types.StringValue(CreateIdentityID(plan.LocalPart, plan.DomainName, plan.Identity))
	plan.Address = custom_types.NewEmailAddressValue(updatedIdentity.Address)
	plan.Name = types.StringValue(updatedIdentity.Name)
	plan.MaySend = types.BoolValue(updatedIdentity.MaySend)
	plan.MayReceive = types.BoolValue(updatedIdentity.MayReceive)
	plan.MayAccessImap = types.BoolValue(updatedIdentity.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(updatedIdentity.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(updatedIdentity.MayAccessManageSieve)
	plan.FooterActive = types.BoolValue(updatedIdentity.FooterActive)
	plan.FooterPlainBody = types.StringValue(updatedIdentity.FooterPlainBody)
	plan.FooterHtmlBody = types.StringValue(updatedIdentity.FooterHtmlBody)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (r *IdentityResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state IdentityResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := r.MigaduClient.DeleteIdentity(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString(), state.Identity.ValueString())
	if err != nil {
		response.Diagnostics.Append(IdentityDeleteError(err))
		return
	}
}

func (r *IdentityResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	idParts := strings.Split(request.ID, "@")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.Append(IdentityImportError(request.ID))
		return
	}

	localPart := idParts[0]
	domainPart := strings.Split(idParts[1], "/")

	if len(domainPart) != 2 || domainPart[0] == "" || domainPart[1] == "" {
		response.Diagnostics.Append(IdentityImportError(request.ID))
		return
	}

	domainName := domainPart[0]
	identity := domainPart[1]

	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"local_part":  localPart,
		"domain_name": domainName,
		"identity":    identity,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("local_part"), localPart)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("domain_name"), domainName)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("identity"), identity)...)
}
