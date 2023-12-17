/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = (*RewriteRuleResource)(nil)
	_ resource.ResourceWithConfigure   = (*RewriteRuleResource)(nil)
	_ resource.ResourceWithImportState = (*RewriteRuleResource)(nil)
)

func NewRewriteRuleResource() resource.Resource {
	return &RewriteRuleResource{}
}

type RewriteRuleResource struct {
	MigaduClient *client.MigaduClient
}

type RewriteRuleResourceModel struct {
	ID            types.String                      `tfsdk:"id"`
	DomainName    custom_types.DomainNameValue      `tfsdk:"domain_name"`
	Name          types.String                      `tfsdk:"name"`
	LocalPartRule types.String                      `tfsdk:"local_part_rule"`
	OrderNum      types.Int64                       `tfsdk:"order_num"`
	Destinations  custom_types.EmailAddressSetValue `tfsdk:"destinations"`
}

func (r *RewriteRuleResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rewrite_rule"
}

func (r *RewriteRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "Provides a rewrite rule.",
		MarkdownDescription: "Provides a rewrite rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Contains the value 'domain_name/name'.",
				MarkdownDescription: "Contains the value `domain_name/name`.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"local_part_rule": schema.StringAttribute{
				Description:         "The regular expression matching the local part of incoming emails.",
				MarkdownDescription: "The regular expression matching the local part of incoming emails.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"domain_name": schema.StringAttribute{
				Description:         "The domain name of the rewrite rule.",
				MarkdownDescription: "The domain name of the rewrite rule.",
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
			"name": schema.StringAttribute{
				Description:         "The name (slug) of the rewrite rule.",
				MarkdownDescription: "The name (slug) of the rewrite rule.",
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
			"order_num": schema.Int64Attribute{
				Description:         "The order of the rewrite rule. Lowest will be executed first.",
				MarkdownDescription: "The order of the rewrite rule. Lowest will be executed first.",
				Required:            false,
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"destinations": schema.SetAttribute{
				Description:         "The destinations of the rewrite rule.",
				MarkdownDescription: "The destinations of the rewrite rule.",
				Required:            true,
				Optional:            false,
				Computed:            false,
				CustomType: custom_types.EmailAddressSetType{
					SetType: types.SetType{
						ElemType: custom_types.EmailAddressType{},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}
func (r *RewriteRuleResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *RewriteRuleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan RewriteRuleResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var destinations []string
	if !plan.Destinations.IsUnknown() {
		response.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	rewrite := &model.RewriteRule{
		Name:          plan.Name.ValueString(),
		LocalPartRule: plan.LocalPartRule.ValueString(),
		OrderNum:      plan.OrderNum.ValueInt64(),
		Destinations:  destinations,
	}

	createdRewrite, err := r.MigaduClient.CreateRewriteRule(ctx, plan.DomainName.ValueString(), rewrite)
	if err != nil {
		response.Diagnostics.Append(RewriteRuleCreateError(err))
		return
	}

	plan.ID = types.StringValue(CreateRewriteRuleID(plan.DomainName, plan.Name))
	plan.Name = types.StringValue(createdRewrite.Name)
	plan.LocalPartRule = types.StringValue(createdRewrite.LocalPartRule)
	plan.OrderNum = types.Int64Value(createdRewrite.OrderNum)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (r *RewriteRuleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state RewriteRuleResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	rewrite, err := r.MigaduClient.GetRewriteRule(ctx, state.DomainName.ValueString(), state.Name.ValueString())
	if err != nil {
		var requestError *client.RequestError
		if errors.As(err, &requestError) {
			if requestError.StatusCode == http.StatusNotFound {
				response.State.RemoveResource(ctx)
				return
			}
		}
		response.Diagnostics.Append(RewriteRuleReadError(err))
		return
	}

	receivedDestinations, diags := custom_types.NewEmailAddressSetValueFrom(ctx, rewrite.Destinations)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if equal, _ := state.Destinations.SetSemanticEquals(ctx, receivedDestinations); !equal {
		state.Destinations = receivedDestinations
	}

	state.ID = types.StringValue(CreateRewriteRuleID(state.DomainName, state.Name))
	state.Name = types.StringValue(rewrite.Name)
	state.LocalPartRule = types.StringValue(rewrite.LocalPartRule)
	state.OrderNum = types.Int64Value(rewrite.OrderNum)

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *RewriteRuleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan RewriteRuleResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var destinations []string
	response.Diagnostics.Append(plan.Destinations.ElementsAs(ctx, &destinations, false)...)
	if response.Diagnostics.HasError() {
		return
	}

	rewrite := &model.RewriteRule{
		Name:          plan.Name.ValueString(),
		LocalPartRule: plan.LocalPartRule.ValueString(),
		OrderNum:      plan.OrderNum.ValueInt64(),
		Destinations:  destinations,
	}

	updatedRewrite, err := r.MigaduClient.UpdateRewriteRule(ctx, plan.DomainName.ValueString(), plan.Name.ValueString(), rewrite)
	if err != nil {
		response.Diagnostics.Append(RewriteRuleUpdateError(err))
		return
	}

	plan.ID = types.StringValue(CreateRewriteRuleID(plan.DomainName, plan.Name))
	plan.Name = types.StringValue(updatedRewrite.Name)
	plan.LocalPartRule = types.StringValue(updatedRewrite.LocalPartRule)
	plan.OrderNum = types.Int64Value(updatedRewrite.OrderNum)

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (r *RewriteRuleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state RewriteRuleResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := r.MigaduClient.DeleteRewriteRule(ctx, state.DomainName.ValueString(), state.Name.ValueString())
	if err != nil {
		response.Diagnostics.Append(RewriteRuleDeleteError(err))
		return
	}
}

func (r *RewriteRuleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	idParts := strings.Split(request.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.Append(RewriteRuleImportError(request.ID))
		return
	}

	domainName := idParts[0]
	name := idParts[1]
	tflog.Trace(ctx, "parsed import ID", map[string]interface{}{
		"domain_name": domainName,
		"name":        name,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("domain_name"), domainName)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("name"), name)...)
}
