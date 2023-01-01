/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                = &rewriteResource{}
	_ resource.ResourceWithConfigure   = &rewriteResource{}
	_ resource.ResourceWithImportState = &rewriteResource{}
)

func NewRewriteResource() resource.Resource {
	return &rewriteResource{}
}

type rewriteResource struct {
	migaduClient *client.MigaduClient
}

type rewriteResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	DomainName           types.String `tfsdk:"domain_name"`
	Name                 types.String `tfsdk:"name"`
	LocalPartRule        types.String `tfsdk:"local_part_rule"`
	OrderNum             types.Int64  `tfsdk:"order_num"`
	Destinations         types.List   `tfsdk:"destinations"`
	DestinationsPunycode types.List   `tfsdk:"destinations_punycode"`
}

func (r *rewriteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rewrite"
}

func (r *rewriteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage a single rewrite rule.",
		MarkdownDescription: "Manage a single rewrite rule.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the rewrite rule to manage.",
				MarkdownDescription: "The domain name of the rewrite rule to manage.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Contains the value 'name'.",
				MarkdownDescription: "Contains the value `name`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name (slug) of the rewrite rule.",
				MarkdownDescription: "The name (slug) of the rewrite rule.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"local_part_rule": schema.StringAttribute{
				Description:         "The regular expression matching the local part of incoming emails.",
				MarkdownDescription: "The regular expression matching the local part of incoming emails.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"order_num": schema.Int64Attribute{
				Description:         "The order of the rewrite rule. Lowest will be executed first.",
				MarkdownDescription: "The order of the rewrite rule. Lowest will be executed first.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"destinations": schema.ListAttribute{
				Description:         "List of email addresses that act as destinations of the rewrite rule.",
				MarkdownDescription: "List of email addresses that act as destinations of the rewrite rule.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ExactlyOneOf(path.MatchRoot("destinations_punycode")),
					listvalidator.SizeAtLeast(1),
				},
			},
			"destinations_punycode": schema.ListAttribute{
				Description:         "List of email addresses that act as destinations of the rewrite rule. Use this attribute instead of 'destinations' in case you want/must use the punycode representation of your domain.",
				MarkdownDescription: "List of email addresses that act as destinations of the rewrite rule. Use this attribute instead of `destinations` in case you want/must use the punycode representation of your domain.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ExactlyOneOf(path.MatchRoot("destinations")),
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}
func (r *rewriteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *rewriteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rewriteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var destinations []string
	if !plan.Destinations.IsUnknown() {
		resp.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.DestinationsPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.DestinationsPunycode.ElementsAs(ctx, &destinations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	rewrite := &model.Rewrite{
		Name:          plan.Name.ValueString(),
		LocalPartRule: plan.LocalPartRule.ValueString(),
		OrderNum:      plan.OrderNum.ValueInt64(),
		Destinations:  destinations,
	}

	createdRewrite, err := r.migaduClient.CreateRewrite(ctx, plan.DomainName.ValueString(), rewrite)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating rewrite rule",
			fmt.Sprintf("Could not create rewrite rule %s: %v", createRewriteID(plan.DomainName, plan.Name), err),
		)
		return
	}

	receivedDestinations, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(createdRewrite.Destinations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	receivedDestinationsPunycode, diags := types.ListValueFrom(ctx, types.StringType, createdRewrite.Destinations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Destinations = receivedDestinations
	plan.DestinationsPunycode = receivedDestinationsPunycode
	plan.ID = types.StringValue(createRewriteID(plan.DomainName, plan.Name))
	plan.Name = types.StringValue(createdRewrite.Name)
	plan.LocalPartRule = types.StringValue(createdRewrite.LocalPartRule)
	plan.OrderNum = types.Int64Value(createdRewrite.OrderNum)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *rewriteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rewriteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rewrite, err := r.migaduClient.GetRewrite(ctx, state.DomainName.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading rewrite rule",
			fmt.Sprintf("Could not read rewrite rule %s: %v", createRewriteID(state.DomainName, state.Name), err),
		)
		return
	}

	receivedDestinations, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(rewrite.Destinations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	receivedDestinationsPunycode, diags := types.ListValueFrom(ctx, types.StringType, rewrite.Destinations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Destinations = receivedDestinations
	state.DestinationsPunycode = receivedDestinationsPunycode
	state.ID = types.StringValue(createRewriteID(state.DomainName, state.Name))
	state.Name = types.StringValue(rewrite.Name)
	state.LocalPartRule = types.StringValue(rewrite.LocalPartRule)
	state.OrderNum = types.Int64Value(rewrite.OrderNum)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *rewriteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rewriteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var destinations []string
	if !plan.Destinations.IsUnknown() {
		resp.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.DestinationsPunycode.IsUnknown() {
		resp.Diagnostics.Append(plan.DestinationsPunycode.ElementsAs(ctx, &destinations, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	rewrite := &model.Rewrite{
		Name:          plan.Name.ValueString(),
		LocalPartRule: plan.LocalPartRule.ValueString(),
		OrderNum:      plan.OrderNum.ValueInt64(),
		Destinations:  destinations,
	}

	updatedRewrite, err := r.migaduClient.UpdateRewrite(ctx, plan.DomainName.ValueString(), plan.Name.ValueString(), rewrite)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating rewrite rule",
			fmt.Sprintf("Could not update rewrite rule %s: %v", createRewriteID(plan.DomainName, plan.Name), err),
		)
		return
	}

	receivedDestinations, diags := types.ListValueFrom(ctx, types.StringType, ConvertEmailsToUnicode(updatedRewrite.Destinations, &resp.Diagnostics))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	receivedDestinationsPunycode, diags := types.ListValueFrom(ctx, types.StringType, updatedRewrite.Destinations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Destinations = receivedDestinations
	plan.DestinationsPunycode = receivedDestinationsPunycode
	plan.ID = types.StringValue(createRewriteID(plan.DomainName, plan.Name))
	plan.Name = types.StringValue(updatedRewrite.Name)
	plan.LocalPartRule = types.StringValue(updatedRewrite.LocalPartRule)
	plan.OrderNum = types.Int64Value(updatedRewrite.OrderNum)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *rewriteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state rewriteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.migaduClient.DeleteRewrite(ctx, state.DomainName.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting rewrite rule",
			fmt.Sprintf("Could not delete rewrite rule %s: %v", createRewriteID(state.DomainName, state.Name), err),
		)
		return
	}
}

func (r *rewriteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Error importing rewrite rule",
			fmt.Sprintf("Expected import identifier with format: 'domain_name/name' Got: '%q'", req.ID),
		)
		return
	}

	domainName := idParts[0]
	name := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"domain_name": domainName,
		"name":        name,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_name"), domainName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
}

func createRewriteID(domainName, name types.String) string {
	return fmt.Sprintf("%s/%s", domainName.ValueString(), name.ValueString())
}
