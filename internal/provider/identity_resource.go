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
	"github.com/metio/terraform-provider-migadu/internal/migadu/client"
	"github.com/metio/terraform-provider-migadu/internal/migadu/model"
	"strings"
)

var (
	_ resource.Resource                = &identityResource{}
	_ resource.ResourceWithConfigure   = &identityResource{}
	_ resource.ResourceWithImportState = &identityResource{}
)

func NewIdentityResource() resource.Resource {
	return &identityResource{}
}

type identityResource struct {
	migaduClient *client.MigaduClient
}

type identityResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	LocalPart            types.String `tfsdk:"local_part"`
	DomainName           types.String `tfsdk:"domain_name"`
	Identity             types.String `tfsdk:"identity"`
	Address              types.String `tfsdk:"address"`
	Name                 types.String `tfsdk:"name"`
	MaySend              types.Bool   `tfsdk:"may_send"`
	MayReceive           types.Bool   `tfsdk:"may_receive"`
	MayAccessImap        types.Bool   `tfsdk:"may_access_imap"`
	MayAccessPop3        types.Bool   `tfsdk:"may_access_pop3"`
	MayAccessManageSieve types.Bool   `tfsdk:"may_access_manage_sieve"`
	Password             types.String `tfsdk:"password"`
	FooterActive         types.Bool   `tfsdk:"footer_active"`
	FooterPlainBody      types.String `tfsdk:"footer_plain_body"`
	FooterHtmlBody       types.String `tfsdk:"footer_html_body"`
}

func (r *identityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity"
}

func (r *identityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage a single identity.",
		MarkdownDescription: "Manage a single identity.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the identity to manage.",
				MarkdownDescription: "The domain name of the identity to manage.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"local_part": schema.StringAttribute{
				Description:         "The local part of the mailbox that owns the identity.",
				MarkdownDescription: "The local part of the mailbox that owns the identity.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identity": schema.StringAttribute{
				Description:         "The local part of the identity to manage.",
				MarkdownDescription: "The local part of the identity to manage.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Contains the value 'local_part@domain_name/identity'.",
				MarkdownDescription: "Contains the value `local_part@domain_name/identity`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				Description:         "Contains the email address of the identity 'identity@domain_name' as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
				MarkdownDescription: "Contains the email address of the identity `identity@domain_name` as returned by the Migadu API. The Migadu API always returns the punycode version of a domain.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
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
func (r *identityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *identityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identityResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	createdIdentity, err := r.migaduClient.CreateIdentity(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), identity)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating identity",
			fmt.Sprintf("Could not create identity %s: %v", createIdentityID(plan.LocalPart, plan.DomainName, plan.Identity), err),
		)
		return
	}

	plan.ID = types.StringValue(createIdentityID(plan.LocalPart, plan.DomainName, plan.Identity))
	plan.Address = types.StringValue(createdIdentity.Address)
	plan.Name = types.StringValue(createdIdentity.Name)
	plan.MaySend = types.BoolValue(createdIdentity.MaySend)
	plan.MayReceive = types.BoolValue(createdIdentity.MayReceive)
	plan.MayAccessImap = types.BoolValue(createdIdentity.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(createdIdentity.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(createdIdentity.MayAccessManageSieve)
	//plan.Password = types.StringValue(createdIdentity.Password)
	plan.FooterActive = types.BoolValue(createdIdentity.FooterActive)
	plan.FooterPlainBody = types.StringValue(createdIdentity.FooterPlainBody)
	plan.FooterHtmlBody = types.StringValue(createdIdentity.FooterHtmlBody)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *identityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity, err := r.migaduClient.GetIdentity(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString(), state.Identity.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading identity",
			fmt.Sprintf("Could not read identity %s: %v", createIdentityID(state.LocalPart, state.DomainName, state.Identity), err),
		)
		return
	}

	state.ID = types.StringValue(createIdentityID(state.LocalPart, state.DomainName, state.Identity))
	state.Address = types.StringValue(identity.Address)
	state.Name = types.StringValue(identity.Name)
	state.MaySend = types.BoolValue(identity.MaySend)
	state.MayReceive = types.BoolValue(identity.MayReceive)
	state.MayAccessImap = types.BoolValue(identity.MayAccessImap)
	state.MayAccessPop3 = types.BoolValue(identity.MayAccessPop3)
	state.MayAccessManageSieve = types.BoolValue(identity.MayAccessManageSieve)
	//state.Password = types.StringValue(identity.Password)
	state.FooterActive = types.BoolValue(identity.FooterActive)
	state.FooterPlainBody = types.StringValue(identity.FooterPlainBody)
	state.FooterHtmlBody = types.StringValue(identity.FooterHtmlBody)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *identityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan identityResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	updatedIdentity, err := r.migaduClient.UpdateIdentity(ctx, plan.DomainName.ValueString(), plan.LocalPart.ValueString(), plan.Identity.ValueString(), identity)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating identity",
			fmt.Sprintf("Could not update identity %s: %v", createIdentityID(plan.LocalPart, plan.DomainName, plan.Identity), err),
		)
		return
	}

	plan.ID = types.StringValue(createIdentityID(plan.LocalPart, plan.DomainName, plan.Identity))
	plan.Address = types.StringValue(updatedIdentity.Address)
	plan.Name = types.StringValue(updatedIdentity.Name)
	plan.MaySend = types.BoolValue(updatedIdentity.MaySend)
	plan.MayReceive = types.BoolValue(updatedIdentity.MayReceive)
	plan.MayAccessImap = types.BoolValue(updatedIdentity.MayAccessImap)
	plan.MayAccessPop3 = types.BoolValue(updatedIdentity.MayAccessPop3)
	plan.MayAccessManageSieve = types.BoolValue(updatedIdentity.MayAccessManageSieve)
	//plan.Password = types.StringValue(updatedIdentity.Password)
	plan.FooterActive = types.BoolValue(updatedIdentity.FooterActive)
	plan.FooterPlainBody = types.StringValue(updatedIdentity.FooterPlainBody)
	plan.FooterHtmlBody = types.StringValue(updatedIdentity.FooterHtmlBody)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *identityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state identityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.migaduClient.DeleteIdentity(ctx, state.DomainName.ValueString(), state.LocalPart.ValueString(), state.Identity.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting identity",
			fmt.Sprintf("Could not delete identity %s: %v", createIdentityID(state.LocalPart, state.DomainName, state.Identity), err),
		)
		return
	}
}

func (r *identityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "@")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Error importing identity",
			fmt.Sprintf("Expected import identifier with format: 'local_part@domain_name/identity' Got: '%q'", req.ID),
		)
		return
	}

	localPart := idParts[0]
	domainName := idParts[1]
	identity := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"local_part":  localPart,
		"domain_name": domainName,
		"identity":    identity,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("local_part"), localPart)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_name"), domainName)...)
}

func createIdentityID(localPart, domainName, identity types.String) string {
	return fmt.Sprintf("%s@%s/%s", localPart.ValueString(), domainName.ValueString(), identity.ValueString())
}
